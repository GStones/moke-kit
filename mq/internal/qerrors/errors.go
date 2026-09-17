package qerrors

import "errors"

var (
	// ErrSubscriptionFailure error encountered subscribing to topic
	ErrSubscriptionFailure = errors.New("ErrSubscriptionFailure")
	// ErrMQTypeUnsupported specified message queue type not supported
	ErrMQTypeUnsupported = errors.New("ErrMQTypeUnsupported")
	// ErrTopicParse topic string was unable to be parsed
	ErrTopicParse = errors.New("ErrTopicParse")
	// ErrNoNatsQueue no NATS MQ implementation was provided during dependency injection
	ErrNoNatsQueue = errors.New("ErrNoNatsQueue")
	// ErrNoLocalQueue no Local MQ implementation was provided during dependency injection
	ErrNoLocalQueue = errors.New("ErrNoLocalQueue")
	// ErrInvalidSubscription invalid Subscription
	ErrInvalidSubscription = errors.New("ErrInvalidSubscription")
	// ErrEmptyTopic empty topic value passed in as argument
	ErrEmptyTopic = errors.New("ErrEmptyTopic")
	// ErrDelayedPublishUnsupported Delayed publishing not supported.
	ErrDelayedPublishUnsupported = errors.New("ErrDelayedPublishUnsupported")
)
