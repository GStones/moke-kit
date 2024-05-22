package nats

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal/message"
	"github.com/gstones/moke-kit/mq/internal/qerrors"
	"github.com/gstones/moke-kit/mq/miface"
)

type Subscription struct {
	topic   string
	mutex   sync.Mutex
	nSub    *nats.Subscription
	handler miface.SubResponseHandler
}

func NewSubscription(
	ctx context.Context,
	topic string,
	conn *nats.Conn,
	deliverySemantics common.DeliverySemantics,
	queue string,
	handler miface.SubResponseHandler,
) (*Subscription, error) {
	sub := &Subscription{
		topic:   topic,
		handler: handler,
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
		}

	case common.AtMostOnce:
		if queue == "" {
			queue = sub.topic
		}

		if nSub, err := conn.QueueSubscribe(sub.topic, queue, sub.handleMessage); err != nil {
			return nil, err
		} else {
			sub.nSub = nSub
		}

	default:
		return nil, errors.New("Unsupported delivery semantics type: " + string(deliverySemantics))
	}
	go func() {
		<-ctx.Done()
		_ = sub.Unsubscribe()
	}()
	return sub, nil
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
	} else {
		mqMsg := message.NewMessage(uid.String(), nMsg.Subject, nMsg.Data, nil)
		s.handler(mqMsg, nil)
	}
}
