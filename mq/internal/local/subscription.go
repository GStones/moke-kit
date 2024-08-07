package local

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/gstones/moke-kit/mq/common"
	message2 "github.com/gstones/moke-kit/mq/internal/message"
	"github.com/gstones/moke-kit/mq/miface"
)

type Subscription struct {
	topic      string
	handler    miface.SubResponseHandler
	subscriber message.Subscriber
}

func CreateSubscription(
	ctx context.Context,
	topic string,
	handler miface.SubResponseHandler,
	subscriber message.Subscriber,
) (*Subscription, error) {
	sub := &Subscription{topic: topic, handler: handler, subscriber: subscriber}
	msgIn, err := subscriber.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range msgIn {
			m := message2.Msg2Message(topic, msg)
			if code := sub.handler(m, nil); code == common.ConsumeNackTransientFailure {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()
	return sub, nil
}

func (s *Subscription) IsValid() bool {
	return s.subscriber != nil
}

func (s *Subscription) Unsubscribe() error {
	if err := s.subscriber.Close(); err != nil {
		return err
	}
	s.subscriber = nil
	return nil
}
