package qerrors

import "github.com/pkg/errors"

var (
	//// mq package scope

	// ErrSubscriptionFailure error encountered subscribing to topic
	ErrSubscriptionFailure = errors.New("ErrSubscriptionFailure")
	// ErrMQTypeUnsupported specified message queue type not supported
	ErrMQTypeUnsupported = errors.New("ErrMQTypeUnsupported")
	// ErrTopicParse topic string was unable to be parsed
	ErrTopicParse = errors.New("ErrTopicParse")
	// ErrNoKafkaQueue no Kafka MQ implementation was provided during dependency injection
	ErrNoKafkaQueue = errors.New("ErrNoKafkaQueue")
	// ErrNoNsqQueue no Nsq MQ implementation was provided during dependency injection
	ErrNoNsqQueue = errors.New("ErrNoNsqQueue")
	// ErrNoNatsQueue no NATS MQ implementation was provided during dependency injection
	ErrNoNatsQueue = errors.New("ErrNoNatsQueue")
	// ErrNoLocalQueue no Local MQ implementation was provided during dependency injection
	ErrNoLocalQueue = errors.New("ErrNoLocalQueue")
	// ErrInvalidSubscription invalid Subscription
	ErrInvalidSubscription = errors.New("ErrInvalidSubscription")
	// ErrDataAlreadySet data payload already set for PubOptions object
	ErrDataAlreadySet = errors.New("ErrDataAlreadySet")
	// ErrEmptyTopic empty topic value passed in as argument
	ErrEmptyTopic = errors.New("ErrEmptyTopic")
	// ErrSemanticsAlreadySet delivery semantics already set
	ErrSemanticsAlreadySet = errors.New("ErrSemanticsAlreadySet")
	// ErrDelayedPublishUnsupported Delayed publishing not supported.
	ErrDelayedPublishUnsupported = errors.New("ErrDelayedPublishUnsupported")
)
