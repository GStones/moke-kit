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
	// ErrInvalidBalancerCode invalid BalancerCode
	ErrInvalidBalancerCode = errors.New("ErrInvalidBalancerCode")
	// ErrInvalidSubscription invalid Subscription
	ErrInvalidSubscription = errors.New("ErrInvalidSubscription")
	// ErrDataAlreadySet data payload already set for PubOptions object
	ErrDataAlreadySet = errors.New("ErrDataAlreadySet")
	// ErrDecoderAlreadySet Decoder already set for SubOptions object
	ErrDecoderAlreadySet = errors.New("ErrDecoderAlreadySet")
	// ErrValuePtrAlreadySet  already set for SubOptions object
	ErrValuePtrAlreadySet = errors.New("ErrValuePtrAlreadySet")
	// ErrEmptyTopic empty topic value passed in as argument
	ErrEmptyTopic = errors.New("ErrEmptyTopic")
	// ErrSemanticsAlreadySet delivery semantics already set
	ErrSemanticsAlreadySet = errors.New("ErrSemanticsAlreadySet")
	// ErrDelayedPublishUnsupported Delayed publishing not supported.
	ErrDelayedPublishUnsupported = errors.New("ErrDelayedPublishUnsupported")
	// ErrAtMostOnceUnsupported AtMostOnce not supported
	ErrAtMostOnceUnsupported = errors.New("ErrAtMostOnceUnsupported")

	//// kafka package scope

	// ErrReaderGroupNotFound no readerGroup was found with the given id
	ErrReaderGroupNotFound = errors.New("ErrReaderGroupNotFound")
	// ErrReaderGroupAlreadyExists a readerGroup already exists for the given id
	ErrReaderGroupAlreadyExists = errors.New("ErrReaderGroupAlreadyExists")
	// ErrUnsupportedConsumeCode handler callback returned an unsupported ConsumptionCode
	ErrUnsupportedConsumeCode = errors.New("ErrUnsupportedConsumeCode")

	//// mock package scope

	// ErrEmptyQueueValue empty queue value passed in for queue subscription
	ErrEmptyQueueValue = errors.New("ErrEmptyQueueValue")
	// ErrSubNotFound subscription not found within sub group
	ErrSubNotFound = errors.New("ErrSubNotFound")
	// ErrUnexpectedRegexResults unexpected regex result array length
	ErrUnexpectedRegexResults = errors.New("ErrUnexpectedRegexResults")
	// ErrTopicValidationFail topic format did not pass regex semantics validation
	ErrTopicValidationFail = errors.New("ErrTopicValidationFail")

	// CLI suite nerrors
	ErrInvalidImplementation    = errors.New("ErrInvalidImplementation")
	ErrInvalidDeliverySemantics = errors.New("ErrInvalidDeliverySemantics")
)
