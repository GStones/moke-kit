package internal

import (
	"strings"

	"github.com/gstones/moke-kit/mq/miface"

	"github.com/pkg/errors"

	"github.com/gstones/moke-kit/mq/internal/qerrors"
)

type MessageQueue struct {
	kafkaMQ miface.MessageQueue
	natsMQ  miface.MessageQueue
	nsqMQ   miface.MessageQueue
	localMQ miface.MessageQueue
}

func NewMessageQueue(
	kafkaMQ miface.MessageQueue,
	natsMQ miface.MessageQueue,
	nsqMQ miface.MessageQueue,
	localMQ miface.MessageQueue,
) *MessageQueue {
	return &MessageQueue{
		kafkaMQ: kafkaMQ,
		natsMQ:  natsMQ,
		nsqMQ:   nsqMQ,
		localMQ: localMQ,
	}
}

func (m *MessageQueue) Subscribe(
	topic string,
	handler miface.SubResponseHandler,
	opts ...miface.SubOption,
) (miface.Subscription, error) {
	if mqType, t, err := parseTopic(topic); err != nil {
		return nil, err
	} else {
		switch mqType {
		case kafka:
			if m.kafkaMQ == nil {
				return nil, qerrors.ErrNoKafkaQueue
			}

			if sub, err := m.kafkaMQ.Subscribe(t, handler, opts...); err != nil {
				return nil, errors.Wrap(err, qerrors.ErrSubscriptionFailure.Error())
			} else {
				return sub, nil
			}

		case nats:
			if m.natsMQ == nil {
				return nil, qerrors.ErrNoNatsQueue
			}

			if sub, err := m.natsMQ.Subscribe(t, handler, opts...); err != nil {
				return nil, errors.Wrap(err, qerrors.ErrSubscriptionFailure.Error())
			} else {
				return sub, nil
			}
		case nsq:
			if m.nsqMQ == nil {
				return nil, qerrors.ErrNoNsqQueue
			}
			if sub, err := m.nsqMQ.Subscribe(t, handler, opts...); err != nil {
				return nil, errors.Wrap(err, qerrors.ErrSubscriptionFailure.Error())
			} else {
				return sub, nil
			}
		case local:
			if m.localMQ == nil {
				return nil, qerrors.ErrNoLocalQueue
			}

			if sub, err := m.localMQ.Subscribe(t, handler, opts...); err != nil {
				return nil, errors.Wrap(err, qerrors.ErrSubscriptionFailure.Error())
			} else {
				return sub, nil
			}

		default:
			return nil, qerrors.ErrMQTypeUnsupported
		}
	}
}

func (m *MessageQueue) Publish(topic string, opts ...miface.PubOption) error {
	if mqType, t, err := parseTopic(topic); err != nil {
		return err
	} else {
		switch mqType {
		case kafka:
			if m.kafkaMQ == nil {
				return qerrors.ErrNoKafkaQueue
			}

			return m.kafkaMQ.Publish(t, opts...)

		case nats:
			if m.natsMQ == nil {
				return qerrors.ErrNoNatsQueue
			}

			return m.natsMQ.Publish(t, opts...)
		case nsq:
			if m.nsqMQ == nil {
				return qerrors.ErrNoNsqQueue
			}
			return m.nsqMQ.Publish(t, opts...)
		case local:
			if m.localMQ == nil {
				return qerrors.ErrNoLocalQueue
			}

			return m.localMQ.Publish(t, opts...)

		default:
			return qerrors.ErrMQTypeUnsupported
		}
	}
}

// topic string should follow the syntax of:
// kafka://topic-name
// nats://some-other-topic
func parseTopic(topic string) (mqType, string, error) {
	sep := "://"

	if len(topic) < 2+len(sep) {
		return unknown, "", qerrors.ErrTopicParse
	} else if elements := strings.Split(topic, sep); len(elements) != 2 {
		return unknown, "", qerrors.ErrTopicParse
	} else if elements[0] == "kafka" {
		return kafka, elements[1], nil
	} else if elements[0] == "nats" {
		return nats, elements[1], nil
	} else if elements[0] == "nsq" {
		return nsq, elements[1], nil
	} else if elements[0] == "local" {
		return local, elements[1], nil
	} else {
		return unknown, elements[1], nil
	}
}

type mqType = int32

const (
	kafka mqType = iota
	nats
	nsq
	local
	unknown
)
