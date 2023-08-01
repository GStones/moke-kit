package internal

import (
	"moke-kit/mq"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type subscription struct {
	topic       string
	mutex       sync.Mutex
	nSub        *nats.Subscription
	handler     mq.SubResponseHandler
	decoder     mq.Decoder
	vPtrFactory mq.ValuePtrFactory
}

func NewSubscription(
	topic string,
	conn *nats.Conn,
	deliverySemantics mq.DeliverySemantics,
	queue string,
	handler mq.SubResponseHandler,
	decoder mq.Decoder,
	vPtrFactory mq.ValuePtrFactory,
) (*subscription, error) {
	sub := &subscription{
		topic:       topic,
		handler:     handler,
		decoder:     decoder,
		vPtrFactory: vPtrFactory,
	}

	if deliverySemantics == mq.Unset {
		deliverySemantics = defaultDeliverySemantics
	}

	switch deliverySemantics {
	case mq.AtLeastOnce:
		if nSub, err := conn.Subscribe(sub.topic, sub.handleMessage); err != nil {
			return nil, err
		} else {
			sub.nSub = nSub
			return sub, nil
		}

	case mq.AtMostOnce:
		if queue == "" {
			queue = sub.topic
		}

		if nSub, err := conn.QueueSubscribe(sub.topic, queue, sub.handleMessage); err != nil {
			return nil, err
		} else {
			sub.nSub = nSub
			return sub, nil
		}

	default:
		return nil, errors.New("Unsupported delivery semantics type: " + string(deliverySemantics))
	}
}

func (s *subscription) IsValid() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.nSub == nil {
		return false
	} else {
		return s.nSub.IsValid()
	}
}

func (s *subscription) Unsubscribe() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.nSub == nil {
		return mq.ErrInvalidSubscription
	} else {
		err := s.nSub.Unsubscribe()
		s.handler = nil

		return err
	}
}

// The NATS mq implementation doesn't care about the mq.ConsumptionCode returned by the handler
func (s *subscription) handleMessage(nMsg *nats.Msg) {
	// Handle the case where Unsubscribe() was called just prior to NATS calling handleMessage()
	if s.handler == nil {
		return
	}

	if uuid, err := utils.NewUUIDString(); err != nil {
		s.handler(nil, err)
	} else if s.decoder != nil && s.vPtrFactory != nil {
		vPtrMessage := s.vPtrFactory.NewVPtr()

		if err := s.decoder.Decode(nMsg.Subject, nMsg.Data, vPtrMessage); err != nil {
			s.handler(nil, err)
		} else {
			mqMsg := NewMessage(uuid, nMsg.Subject, nil, vPtrMessage)
			s.handler(mqMsg, nil)
		}
	} else {
		mqMsg := NewMessage(uuid, nMsg.Subject, nMsg.Data, nil)

		s.handler(mqMsg, nil)
	}
}
