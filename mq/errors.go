package mq

import "github.com/pkg/errors"

var (
	//// mq package scope

	// error encountered subscribing to topic
	ErrSubscriptionFailure = errors.New("ErrSubscriptionFailure")
	// specified message queue type not supported
	ErrMQTypeUnsupported = errors.New("ErrMQTypeUnsupported")
	// topic string was unable to be parsed
	ErrTopicParse = errors.New("ErrTopicParse")
	// no Kafka MQ implementation was provided during dependency injection
	ErrNoKafkaQueue = errors.New("ErrNoKafkaQueue")
	// no Nsq MQ implementation was provided during dependency injection
	ErrNoNsqQueue = errors.New("ErrNoNsqQueue")
	// no NATS MQ implementation was provided during dependency injection
	ErrNoNatsQueue = errors.New("ErrNoNatsQueue")
	// no Local MQ implementation was provided during dependency injection
	ErrNoLocalQueue = errors.New("ErrNoLocalQueue")
	// invalid BalancerCode
	ErrInvalidBalancerCode = errors.New("ErrInvalidBalancerCode")
	// invalid Subscription
	ErrInvalidSubscription = errors.New("ErrInvalidSubscription")
	// data payload already set for PubOptions object
	ErrDataAlreadySet = errors.New("ErrDataAlreadySet")
	// Decoder already set for SubOptions object
	ErrDecoderAlreadySet = errors.New("ErrDecoderAlreadySet")
	// ValuePtrFactory already set for SubOptions object
	ErrValuePtrAlreadySet = errors.New("ErrValuePtrAlreadySet")
	// "empty topic value passed in as argument
	ErrEmptyTopic = errors.New("ErrEmptyTopic")
	// delivery semantics already set
	ErrSemanticsAlreadySet = errors.New("ErrSemanticsAlreadySet")
	// Delayed publishing not supported.
	ErrDelayedPublishUnsupported = errors.New("ErrDelayedPublishUnsupported")
	// AtMostOnce not supported
	ErrAtMostOnceUnsupported = errors.New("ErrAtMostOnceUnsupported")

	//// kafka package scope

	// no readerGroup was found with the given id
	ErrReaderGroupNotFound = errors.New("ErrReaderGroupNotFound")
	// a readerGroup already exists for the given id
	ErrReaderGroupAlreadyExists = errors.New("ErrReaderGroupAlreadyExists")
	// handler callback returned an unsupported ConsumptionCode
	ErrUnsupportedConsumeCode = errors.New("ErrUnsupportedConsumeCode")

	//// mock package scope

	// empty queue value passed in for queue subscription
	ErrEmptyQueueValue = errors.New("ErrEmptyQueueValue")
	// subscription not found within sub group
	ErrSubNotFound = errors.New("ErrSubNotFound")
	// unexpected regex result array length
	ErrUnexpectedRegexResults = errors.New("ErrUnexpectedRegexResults")
	// topic format did not pass regex semantics validation
	ErrTopicValidationFail = errors.New("ErrTopicValidationFail")

	// CLI suite errors
	ErrInvalidImplementation    = errors.New("ErrInvalidImplementation")
	ErrInvalidDeliverySemantics = errors.New("ErrInvalidDeliverySemantics")
)
