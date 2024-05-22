package nats

import (
	"context"
	"net/url"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal/qerrors"
	"github.com/gstones/moke-kit/mq/miface"
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
	ctx context.Context,
	topic string,
	handler miface.SubResponseHandler,
	sOpts ...miface.SubOption,
) (miface.Subscription, error) {
	if topic == "" {
		return nil, qerrors.ErrEmptyTopic
	} else {
		topic = common.NamespaceTopic(topic)
	}

	if options, err := miface.NewSubOptions(sOpts...); err != nil {
		return nil, err
	} else {
		return NewSubscription(
			ctx,
			topic,
			m.conn,
			options.DeliverySemantics,
			options.GroupId,
			handler,
		)
	}
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
		return m.conn.Publish(topic, options.Data)
	}
}
