package nats

import (
	"context"
	"net/url"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	nc "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal/logger"
	message2 "github.com/gstones/moke-kit/mq/internal/message"
	"github.com/gstones/moke-kit/mq/internal/qerrors"
	"github.com/gstones/moke-kit/mq/miface"
)

type MessageQueue struct {
	logger    *zap.Logger
	subscribe message.Subscriber
	publisher message.Publisher
}

func NewMessageQueue(logger *zap.Logger, address string) (*MessageQueue, error) {
	if u, err := url.Parse(address); err != nil {
		return nil, err
	} else if conn, err := nats.Connect(
		u.String(),
		nats.RetryOnFailedConnect(true),
		nats.ReconnectWait(1*time.Second),
		nats.Timeout(30*time.Second),
	); err != nil {
		return nil, err
	} else {
		mq := &MessageQueue{logger: logger}
		if err := mq.newSubscribe(conn); err != nil {
			return nil, err
		}
		if err := mq.newPublisher(conn); err != nil {
			return nil, err
		}
		return mq, nil
	}
}

func (m *MessageQueue) newSubscribe(conn *nats.Conn) error {
	subscriber, err := nc.NewSubscriberWithNatsConn(
		conn,
		nc.SubscriberSubscriptionConfig{
			CloseTimeout:   30 * time.Second,
			AckWaitTimeout: 30 * time.Second,
			Unmarshaler:    marshaler,
			JetStream:      jsConfig,
		},
		logger.NewZapLoggerAdapter(m.logger),
	)
	if err != nil {
		return err
	}
	m.subscribe = subscriber
	return nil
}

func (m *MessageQueue) newPublisher(conn *nats.Conn) error {
	publisher, err := nc.NewPublisherWithNatsConn(
		conn,
		nc.PublisherPublishConfig{
			Marshaler: marshaler,
			JetStream: jsConfig,
		},
		logger.NewZapLoggerAdapter(m.logger),
	)
	if err != nil {
		return err
	}
	m.publisher = publisher
	return nil
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

	msgChan, err := m.subscribe.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}

	go func() {
		for msg := range msgChan {
			ms := message2.Msg2Message(topic, msg)
			if code := handler(ctx, ms, nil); code == common.ConsumeNackTransientFailure {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()
	return nil, nil
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
