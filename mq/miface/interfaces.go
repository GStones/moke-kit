package miface

import "context"

type MessageQueue interface {
	Subscribe(context context.Context, topic string, handler SubResponseHandler, opts ...SubOption) (Subscription, error)
	Publish(topic string, opts ...PubOption) error
}

type Subscription interface {
	IsValid() bool
	Unsubscribe() error
}

type Message interface {
	ID() string
	Topic() string
	Data() []byte
	VPtr() (vPtr any)
}
