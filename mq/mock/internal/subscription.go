package internal

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/gstones/platform/services/common/mq"
	m "github.com/gstones/platform/services/common/mq/mock/internal/mockmq"
	"github.com/gstones/platform/services/common/utils"
)

type subscription struct {
	topic       string
	mutex       sync.Mutex
	mockSub     m.Subscription
	handler     mq.SubResponseHandler
	decoder     mq.Decoder
	vPtrFactory mq.ValuePtrFactory
}

func NewSubscription(
	topic string,
	mockMQ m.MessageQueue,
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
		if mockSub, err := mockMQ.Subscribe(sub.topic, sub.handleMessage); err != nil {
			return nil, err
		} else {
			sub.mockSub = mockSub
			return sub, nil
		}

	case mq.AtMostOnce:
		if queue == "" {
			queue = topic
		}

		if mockSub, err := mockMQ.QueueSubscribe(sub.topic, queue, sub.handleMessage); err != nil {
			return nil, err
		} else {
			sub.mockSub = mockSub
			return sub, nil
		}

	default:
		return nil, errors.New("Unsupported delivery semantics type: " + string(deliverySemantics))
	}
}

func (s *subscription) IsValid() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.mockSub == nil {
		return false
	} else {
		return s.mockSub.IsValid()
	}
}

func (s *subscription) Unsubscribe() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.mockSub == nil {
		return mq.ErrInvalidSubscription
	} else {
		err := s.mockSub.Unsubscribe()
		s.handler = nil

		return err
	}
}

// The mock mq implementation doesn't care about the mq.ConsumptionCode returned by the handler
func (s *subscription) handleMessage(message m.Message) {
	// Handle the case where Unsubscribe() was called just prior to the mock calling handleMessage()
	if s.handler == nil {
		return
	}

	if uuid, err := utils.NewUUIDString(); err != nil {
		s.handler(nil, err)
	} else if s.decoder != nil && s.vPtrFactory != nil {
		vPtrMessage := s.vPtrFactory.NewVPtr()

		if err := s.decoder.Decode(message.Topic, message.Data, vPtrMessage); err != nil {
			s.handler(nil, err)
		} else {
			mqMsg := NewMessage(uuid, message.Topic, nil, vPtrMessage)
			s.handler(mqMsg, nil)
		}
	} else {
		mqMsg := NewMessage(uuid, message.Topic, message.Data, nil)

		s.handler(mqMsg, nil)
	}
}
