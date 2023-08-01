package mockmq

import "time"

type MessageQueue interface {
	Subscribe(topic string, handler MessageHandler) (Subscription, error)
	QueueSubscribe(topic string, queue string, handler MessageHandler) (Subscription, error)
	Publish(topic string, data []byte) error
	PublishWithDelay(topic string, data []byte, delay time.Duration) error
}

type Subscription interface {
	IsValid() bool
	Unsubscribe() error
	Queue() string
	SendMessage(message Message) error
}

type Unsubscriber interface {
	RemoveSubscription(sub *subscription) error
}

type Sender interface {
	SendMessage(topic string, data []byte)
}

type PermAtLeastOnceSubGroup interface {
	Unsubscriber
	Sender
	NewSubscription(queue string, handler MessageHandler) Subscription
	Subscriptions() []Subscription
}

type PermAtMostOnceSubGroup interface {
	Unsubscriber
	NewSubscription(queue string, handler MessageHandler) Subscription
	Subscriptions() []Subscription
}

type TempAtMostOnceSubGroup interface {
	Sender
	AddSubscriptions(subs ...Subscription)
}
