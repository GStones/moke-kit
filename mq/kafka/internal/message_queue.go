package internal

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"

	k "github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/mq"
	"github.com/gstones/platform/services/common/network"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type messageQueue struct {
	logger             *zap.Logger
	brokers            Brokers
	numPartitions      int
	replicationFactor  int
	readerStatsPeriod  StatsPeriod
	writersStatsPeriod StatsPeriod

	activeAtLeastOnceRGs map[string]*readerGroup
	activeAtMostOnceIdGs map[string]*idGroup
	argMutex             sync.Mutex

	activeWriters map[string]*writer
	awMutex       sync.Mutex

	deadLetterWriter *deadLetterWriter
}

func NewMessageQueue(
	logger *zap.Logger,
	brokerUrls []string,
) (*messageQueue, error) {
	var brokers []network.Address

	for _, b := range brokerUrls {
		if u, err := url.Parse(b); err != nil {
			return nil, err
		} else if port, err := strconv.Atoi(u.Port()); err != nil {
			return nil, err
		} else {
			address := network.Address{
				Network: network.TCP,
				Host:    u.Hostname(),
				Port:    port,
			}

			brokers = append(brokers, address)
		}
	}

	if dlWriter, err := NewDeadLetterWriter(brokers); err != nil {
		return nil, err
	} else {
		return &messageQueue{
			logger:               logger,
			brokers:              brokers,
			numPartitions:        numPartitions,
			replicationFactor:    replicationFactor,
			readerStatsPeriod:    readerStatsPeriod,
			writersStatsPeriod:   writerStatsPeriod,
			activeAtLeastOnceRGs: make(map[string]*readerGroup),
			activeAtMostOnceIdGs: make(map[string]*idGroup),
			activeWriters:        make(map[string]*writer),
			deadLetterWriter:     dlWriter,
		}, nil
	}
}

// Subscribe does not block caller. The handler will be called with reader results
func (m *messageQueue) Subscribe(
	topic string,
	handler mq.SubResponseHandler,
	opts ...mq.SubOption,
) (mq.Subscription, error) {
	if topic == "" {
		return nil, mq.ErrEmptyTopic
	} else {
		topic = mq.NamespaceTopic(topic)
	}

	if options, err := mq.NewSubOptions(opts...); err != nil {
		return nil, err
	} else if rg, err := m.openReaderGroup(topic, options.DeliverySemantics, options.GroupId); err != nil {
		return nil, err
	} else {
		return rg.Subscribe(handler, options.Decoder, options.VPtrFactory)
	}
}

// Send blocks caller until: The message has been written, or the maximum number of
// attempts was reached, or the context timeout deadline was exceeded.
func (m *messageQueue) Publish(topic string, opts ...mq.PubOption) error {
	if topic == "" {
		return mq.ErrEmptyTopic
	} else {
		topic = mq.NamespaceTopic(topic)
	}

	if options, err := mq.NewPubOptions(opts...); err != nil {
		return err
	} else if options.Delay != 0 {
		return mq.ErrDelayedPublishUnsupported
	} else if writer, err := m.openWriter(topic); err != nil {
		return err
	} else {
		return writer.Publish(topic, options.Data)
	}
}

func (m *messageQueue) openReaderGroup(
	topic string,
	deliverySemantics mq.DeliverySemantics,
	groupId string,
) (readerGroup *readerGroup, err error) {
	if deliverySemantics == mq.Unset {
		deliverySemantics = defaultDeliverySemantics
	}

	m.argMutex.Lock()
	defer m.argMutex.Unlock()

	if groupId, err = makeGroupID(topic, groupId, deliverySemantics); err != nil {
		return nil, err
	}

	switch deliverySemantics {
	case mq.AtLeastOnce:
		if rg, ok := m.activeAtLeastOnceRGs[topic]; ok {
			return rg, nil
		}

	case mq.AtMostOnce:
		if _, ok := m.activeAtMostOnceIdGs[topic]; ok {
			if readerGroup, err = m.activeAtMostOnceIdGs[topic].OpenReaderGroup(groupId); err == nil {
				return readerGroup, nil
			} else if err != nil && err != mq.ErrReaderGroupNotFound {
				return nil, err
			}
		}

	default:
		return nil, errors.New("Unsupported delivery semantics type: " + string(deliverySemantics))
	}

	if err := m.createTopic(topic); err != nil {
		return nil, err
	}

	// If a readerGroup for a given topic + semantics ( + group id, possibly ) combination does not exist, we create it
	readerGroup, err = NewReaderGroup(
		m.logger,
		deliverySemantics,
		groupId,
		topic,
		m.brokers,
		int(m.numPartitions),
		m.readerStatsPeriod,
		m.deadLetterWriter,
	)

	switch deliverySemantics {
	case mq.AtLeastOnce:
		m.activeAtLeastOnceRGs[topic] = readerGroup

	case mq.AtMostOnce:
		if _, ok := m.activeAtMostOnceIdGs[topic]; !ok {
			m.activeAtMostOnceIdGs[topic] = NewIdGroup()
		}

		if err := m.activeAtMostOnceIdGs[topic].AddReaderGroup(groupId, readerGroup); err != nil {
			return readerGroup, err
		}

	default:
		return nil, errors.New("Unsupported delivery semantics type: " + string(deliverySemantics))
	}

	return readerGroup, err
}

func (m *messageQueue) openWriter(topic string) (*writer, error) {
	var openedWriter *writer

	m.awMutex.Lock()
	defer m.awMutex.Unlock()

	if w, ok := m.activeWriters[topic]; ok {
		openedWriter = w
	} else {
		if err := m.createTopic(topic); err != nil {
			return nil, err
		}

		var b Balancer

		switch balancer {
		case BalancerCodeRoundRobin:
			b = &k.RoundRobin{}

		case BalancerCodeLeastBytes:
			b = &k.LeastBytes{}

		case BalancerCodeHash:
			b = &k.Hash{}

		default:
			return nil, mq.ErrInvalidBalancerCode
		}

		writerConfig := k.WriterConfig{
			Brokers:  m.brokers.Strings(),
			Topic:    topic,
			Balancer: b,
		}

		openedWriter = NewWriter(m.logger, writerConfig, m.writersStatsPeriod)

		m.activeWriters[topic] = openedWriter
	}

	return openedWriter, nil
}

func (m *messageQueue) createTopic(topic string) error {
	if _, ok := m.activeWriters[topic]; ok {
		return nil
	}

	if _, ok := m.activeAtLeastOnceRGs[topic]; ok {
		return nil
	}

	if _, ok := m.activeAtMostOnceIdGs[topic]; ok {
		return nil
	}

	// We attempt to create the Topic if we have not yet come across it since service runtime
	topicConfig := k.TopicConfig{
		Topic:              topic,
		NumPartitions:      int(m.numPartitions),
		ReplicationFactor:  int(m.replicationFactor),
		ReplicaAssignments: nil,
		ConfigEntries:      nil,
	}

	broker := m.brokers.Random()

	if kConn, err := k.Dial(broker.Network.String(), broker.String()); err != nil {
		return err
	} else if err = kConn.CreateTopics(topicConfig); err != nil {
		switch err {
		case k.NotController:
			// Ok. In practice this error is received when topics already exist.
			return nil
		default:
			return err
		}
	} else {
		return nil
	}
}

func makeGroupID(topic string, groupId string, semantics mq.DeliverySemantics) (string, error) {
	switch semantics {
	case mq.AtLeastOnce:
		// Assign each consumer a "unique" group ID; each message will be
		// delivered to every consumer
		return fmt.Sprintf("%s-%s", topic, xid.New().String()), nil
	case mq.AtMostOnce:
		// assign each consumer the same group ID; each message will be delivered
		// only to a single consumer
		if groupId != "" {
			return groupId, nil
		} else {
			return topic, nil
		}
	default:
		return "", errors.New("Unsupported delivery semantics type: " + string(semantics))
	}
}
