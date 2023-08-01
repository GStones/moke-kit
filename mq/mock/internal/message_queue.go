package internal

import (
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/mq"
	m "github.com/gstones/platform/services/common/mq/mock/internal/mockmq"
)

type messageQueue struct {
	mockMQ m.MessageQueue
}

func NewMessageQueue(logger *zap.Logger, deployment string) (*messageQueue, error) {
	return &messageQueue{mockMQ: m.NewMessageQueue(logger, deployment)}, nil
}

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
	} else {
		return NewSubscription(
			topic,
			m.mockMQ,
			options.DeliverySemantics,
			options.GroupId,
			handler,
			options.Decoder,
			options.VPtrFactory,
		)
	}
}

func (m *messageQueue) Publish(topic string, opts ...mq.PubOption) error {
	if topic == "" {
		return mq.ErrEmptyTopic
	} else {
		topic = mq.NamespaceTopic(topic)
	}

	if options, err := mq.NewPubOptions(opts...); err != nil {
		return err
	} else {
		if options.Delay == 0 {
			return m.mockMQ.Publish(topic, options.Data)
		} else {
			return m.mockMQ.PublishWithDelay(topic, options.Data, options.Delay)
		}
	}
}
