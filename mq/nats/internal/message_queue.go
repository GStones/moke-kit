package internal

import (
	"net/url"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"moke-kit/mq"
)

type messageQueue struct {
	logger *zap.Logger
	conn   *nats.Conn
}

func NewMessageQueue(logger *zap.Logger, address string) (*messageQueue, error) {
	if u, err := url.Parse(address); err != nil {
		return nil, err
	} else if conn, err := nats.Connect(u.String()); err != nil {
		return nil, err
	} else {
		return &messageQueue{logger: logger, conn: conn}, nil
	}
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
			m.conn,
			options.DeliverySemantics,
			options.GroupId,
			handler,
			options.Decoder,
			options.VPtrFactory)
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
	} else if options.Delay != 0 {
		return mq.ErrDelayedPublishUnsupported
	} else {
		return m.conn.Publish(topic, options.Data)
	}
}
