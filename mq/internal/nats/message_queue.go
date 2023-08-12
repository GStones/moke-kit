package nats

import (
	"net/url"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"moke-kit/mq/common"
	"moke-kit/mq/internal/qerrors"
	"moke-kit/mq/qiface"
)

type MessageQueue struct {
	logger *zap.Logger
	conn   *nats.Conn
}

func NewMessageQueue(logger *zap.Logger, address string) (*MessageQueue, error) {
	if u, err := url.Parse(address); err != nil {
		return nil, err
	} else if conn, err := nats.Connect(u.String()); err != nil {
		return nil, err
	} else {
		return &MessageQueue{logger: logger, conn: conn}, nil
	}
}

func (m *MessageQueue) Subscribe(
	topic string,
	handler qiface.SubResponseHandler,
	sOpts ...qiface.SubOption,
) (qiface.Subscription, error) {
	if topic == "" {
		return nil, qerrors.ErrEmptyTopic
	} else {
		topic = common.NamespaceTopic(topic)
	}

	if options, err := qiface.NewSubOptions(sOpts...); err != nil {
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

func (m *MessageQueue) Publish(topic string, pOpts ...qiface.PubOption) error {
	if topic == "" {
		return qerrors.ErrEmptyTopic
	} else {
		topic = common.NamespaceTopic(topic)
	}

	if options, err := qiface.NewPubOptions(pOpts...); err != nil {
		return err
	} else if options.Delay != 0 {
		return qerrors.ErrDelayedPublishUnsupported
	} else {
		return m.conn.Publish(topic, options.Data)
	}
}
