package local

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/gstones/moke-kit/mq/common"
	message2 "github.com/gstones/moke-kit/mq/internal/message"
	"github.com/gstones/moke-kit/mq/internal/subscription"
	"github.com/gstones/moke-kit/mq/miface"
)

func CreateSubscription(
	ctx context.Context,
	topic string,
	handler miface.SubResponseHandler,
	subscriber message.Subscriber,
) (miface.Subscription, error) {
	subCtx, cancel := context.WithCancel(ctx)
	msgIn, err := subscriber.Subscribe(subCtx, topic)
	if err != nil {
		cancel()
		return nil, err
	}
	sub := subscription.NewCancelSubscription(cancel)
	go func() {
		defer sub.Unsubscribe()
		for msg := range msgIn {
			m := message2.Msg2Message(topic, msg)
			if code := handler(m, nil); code == common.ConsumeNackTransientFailure {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()
	return sub, nil
}
