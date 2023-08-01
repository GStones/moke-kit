package mq

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	KafkaProtocol = "kafka://"
	NatsProtocol  = "nats://"
	NsqProtocol   = "nsq://"
	LocalProtocol = "local://"
)

type messageQueue struct {
	kafkaMQ MessageQueue
	natsMQ  MessageQueue
	nsqMQ   MessageQueue
	localMQ MessageQueue
}

func NewMessageQueue(
	kafkaMQ MessageQueue,
	natsMQ MessageQueue,
	nsqMQ MessageQueue,
	localMQ MessageQueue,
) *messageQueue {
	return &messageQueue{
		kafkaMQ: kafkaMQ,
		natsMQ:  natsMQ,
		nsqMQ:   nsqMQ,
		localMQ: localMQ,
	}
}

func (m *messageQueue) Subscribe(
	topic string,
	handler SubResponseHandler,
	opts ...SubOption,
) (Subscription, error) {
	if mqType, t, err := parseTopic(topic); err != nil {
		return nil, err
	} else {
		switch mqType {
		case kafka:
			if m.kafkaMQ == nil {
				return nil, ErrNoKafkaQueue
			}

			if sub, err := m.kafkaMQ.Subscribe(t, handler, opts...); err != nil {
				return nil, errors.Wrap(err, ErrSubscriptionFailure.Error())
			} else {
				return sub, nil
			}

		case nats:
			if m.natsMQ == nil {
				return nil, ErrNoNatsQueue
			}

			if sub, err := m.natsMQ.Subscribe(t, handler, opts...); err != nil {
				return nil, errors.Wrap(err, ErrSubscriptionFailure.Error())
			} else {
				return sub, nil
			}
		case nsq:
			if m.nsqMQ == nil {
				return nil, ErrNoNsqQueue
			}
			if sub, err := m.nsqMQ.Subscribe(t, handler, opts...); err != nil {
				return nil, errors.Wrap(err, ErrSubscriptionFailure.Error())
			} else {
				return sub, nil
			}
		case local:
			if m.localMQ == nil {
				return nil, ErrNoLocalQueue
			}

			if sub, err := m.localMQ.Subscribe(t, handler, opts...); err != nil {
				return nil, errors.Wrap(err, ErrSubscriptionFailure.Error())
			} else {
				return sub, nil
			}

		default:
			return nil, ErrMQTypeUnsupported
		}
	}
}

func (m *messageQueue) Publish(topic string, opts ...PubOption) error {
	if mqType, t, err := parseTopic(topic); err != nil {
		return err
	} else {
		switch mqType {
		case kafka:
			if m.kafkaMQ == nil {
				return ErrNoKafkaQueue
			}

			return m.kafkaMQ.Publish(t, opts...)

		case nats:
			if m.natsMQ == nil {
				return ErrNoNatsQueue
			}

			return m.natsMQ.Publish(t, opts...)
		case nsq:
			if m.nsqMQ == nil {
				return ErrNoNsqQueue
			}
			return m.nsqMQ.Publish(t, opts...)
		case local:
			if m.localMQ == nil {
				return ErrNoLocalQueue
			}

			return m.localMQ.Publish(t, opts...)

		default:
			return ErrMQTypeUnsupported
		}
	}
}

// topic string should follow the syntax of:
// kafka://topic-name
// nats://some-other-topic
func parseTopic(topic string) (mqType, string, error) {
	sep := "://"

	if len(topic) < 2+len(sep) {
		return unknown, "", ErrTopicParse
	} else if elements := strings.Split(topic, sep); len(elements) != 2 {
		return unknown, "", ErrTopicParse
	} else if elements[0] == "kafka" {
		return kafka, elements[1], nil
	} else if elements[0] == "nats" {
		return nats, elements[1], nil
	} else if elements[0] == "nsq" {
		return nsq, elements[1], nil
	} else if elements[0] == "local" {
		return local, elements[1], nil
	} else {
		return unknown, elements[1], nil
	}
}

type mqType = int32

const (
	kafka mqType = iota
	nats
	nsq
	local
	unknown
)
