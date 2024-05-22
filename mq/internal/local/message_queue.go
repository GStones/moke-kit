package local

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal/qerrors"
	"github.com/gstones/moke-kit/mq/miface"
)

type MessageQueue struct {
	logger     *zap.Logger
	publisher  message.Publisher
	subscriber message.Subscriber
}

func NewMessageQueue(logger *zap.Logger, bufferSize int64, persistent bool, isBlocked bool) *MessageQueue {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{
			OutputChannelBuffer:            bufferSize,
			Persistent:                     persistent,
			BlockPublishUntilSubscriberAck: isBlocked,
		},
		watermill.NewStdLogger(true, true),
	)
	return &MessageQueue{
		logger:     logger,
		publisher:  pubSub,
		subscriber: pubSub,
	}
}

func (m *MessageQueue) Subscribe(
	ctx context.Context,
	topic string,
	handler miface.SubResponseHandler,
	_ ...miface.SubOption,
) (miface.Subscription, error) {
	if topic == "" {
		return nil, qerrors.ErrEmptyTopic
	}
	topic = common.NamespaceTopic(topic)
	return CreateSubscription(ctx, topic, handler, m.subscriber)
}

func (m *MessageQueue) Publish(topic string, pOpts ...miface.PubOption) error {
	if topic == "" {
		return qerrors.ErrEmptyTopic
	} else {
		topic = common.NamespaceTopic(topic)
	}

	if options, err := miface.NewPubOptions(pOpts...); err != nil {
		return err
	} else if options.Delay != 0 {
		return qerrors.ErrDelayedPublishUnsupported
	} else {
		msg := message.NewMessage(watermill.NewUUID(), options.Data)
		return m.publisher.Publish(topic, msg)
	}

}
