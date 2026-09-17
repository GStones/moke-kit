package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/gstones/moke-kit/mq/internal/qerrors"
	"github.com/gstones/moke-kit/mq/miface"
)

type MessageQueue struct {
	natsMQ  miface.MessageQueue
	localMQ miface.MessageQueue
}

func NewMessageQueue(natsMQ, localMQ miface.MessageQueue) *MessageQueue {
	return &MessageQueue{
		natsMQ:  natsMQ,
		localMQ: localMQ,
	}
}

func (m *MessageQueue) Subscribe(
	ctx context.Context,
	topic string,
	handler miface.SubResponseHandler,
	opts ...miface.SubOption,
) (miface.Subscription, error) {
	q, t, err := m.backend(topic)
	if err != nil {
		return nil, err
	}
	sub, err := q.Subscribe(ctx, t, handler, opts...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", qerrors.ErrSubscriptionFailure, err)
	}
	return sub, nil
}

func (m *MessageQueue) Publish(topic string, opts ...miface.PubOption) error {
	q, t, err := m.backend(topic)
	if err != nil {
		return err
	}
	return q.Publish(t, opts...)
}

func (m *MessageQueue) backend(topic string) (miface.MessageQueue, string, error) {
	mqType, t, err := parseTopic(topic)
	if err != nil {
		return nil, "", err
	}
	switch mqType {
	case nats:
		if m.natsMQ == nil {
			return nil, "", qerrors.ErrNoNatsQueue
		}
		return m.natsMQ, t, nil
	case local:
		if m.localMQ == nil {
			return nil, "", qerrors.ErrNoLocalQueue
		}
		return m.localMQ, t, nil
	default:
		return nil, "", qerrors.ErrMQTypeUnsupported
	}
}

// topic string should follow the syntax of:
// nats://some-topic
// local://some-other-topic
func parseTopic(topic string) (mqType, string, error) {
	scheme, name, ok := strings.Cut(topic, "://")
	if !ok || scheme == "" || name == "" {
		return unknown, "", qerrors.ErrTopicParse
	}
	switch scheme {
	case "nats":
		return nats, name, nil
	case "local":
		return local, name, nil
	default:
		return unknown, "", qerrors.ErrMQTypeUnsupported
	}
}

type mqType = int32

const (
	nats mqType = iota
	local
	unknown
)
