package nats

import (
	"sync"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal/qerrors"
	"github.com/gstones/moke-kit/mq/logic"
)

type Subscription struct {
	topic       string
	mutex       sync.Mutex
	nSub        *nats.Subscription
	handler     logic.SubResponseHandler
	decoder     logic.Decoder
	vPtrFactory logic.ValuePtrFactory
}

func NewSubscription(
	topic string,
	conn *nats.Conn,
	deliverySemantics common.DeliverySemantics,
	queue string,
	handler logic.SubResponseHandler,
	decoder logic.Decoder,
	vPtrFactory logic.ValuePtrFactory,
) (*Subscription, error) {
	sub := &Subscription{
		topic:       topic,
		handler:     handler,
		decoder:     decoder,
		vPtrFactory: vPtrFactory,
	}

	if deliverySemantics == common.Unset {
		deliverySemantics = defaultDeliverySemantics
	}

	switch deliverySemantics {
	case common.AtLeastOnce:
		if nSub, err := conn.Subscribe(sub.topic, sub.handleMessage); err != nil {
			return nil, err
		} else {
			sub.nSub = nSub
			return sub, nil
		}

	case common.AtMostOnce:
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

func (s *Subscription) IsValid() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.nSub == nil {
		return false
	} else {
		return s.nSub.IsValid()
	}
}

func (s *Subscription) Unsubscribe() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.nSub == nil {
		return qerrors.ErrInvalidSubscription
	} else {
		err := s.nSub.Unsubscribe()
		s.handler = nil

		return err
	}
}

// The NATS mq implementation doesn't care about the mq.ConsumptionCode returned by the handler
func (s *Subscription) handleMessage(nMsg *nats.Msg) {
	// Handle the case where Unsubscribe() was called just prior to NATS calling handleMessage()
	if s.handler == nil {
		return
	}
	if uid, err := uuid.NewUUID(); err != nil {
		s.handler(nil, err)
	} else if s.decoder != nil && s.vPtrFactory != nil {
		vPtrMessage := s.vPtrFactory.NewVPtr()

		if err := s.decoder.Decode(nMsg.Subject, nMsg.Data, vPtrMessage); err != nil {
			s.handler(nil, err)
		} else {
			mqMsg := NewMessage(uid.String(), nMsg.Subject, nil, vPtrMessage)
			s.handler(mqMsg, nil)
		}
	} else {
		mqMsg := NewMessage(uid.String(), nMsg.Subject, nMsg.Data, nil)

		s.handler(mqMsg, nil)
	}
}
