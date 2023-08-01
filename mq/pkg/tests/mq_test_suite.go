package tests

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/logging"
	"github.com/gstones/platform/services/common/mq"
	"github.com/gstones/platform/services/common/mq/kafka"
	"github.com/gstones/platform/services/common/mq/mock"
	"github.com/gstones/platform/services/common/mq/nats"
	"github.com/gstones/platform/services/common/mq/nsq"
)

var msgReceivedCount int

// Mutex for thread-safe handler responses
// Each handler response increments the message received count by one
var msgReceivedMutex sync.Mutex

const (
	atMostOnceKafkaMsgExpectedCt  = 8
	atLeastOnceKafkaMsgExpectedCt = 8

	atMostOnceLocalMsgExpectedCt  = 8
	atLeastOnceLocalMsgExpectedCt = 33

	atMostOnceNatsMsgExpectedCt  = 8
	atLeastOnceNatsMsgExpectedCt = 33

	atMostOnceNsqMsgExpectedCt  = 8
	atLeastOnceNsqMsgExpectedCt = 8

	defaultKafkaUrl       = "tcp://100.69.89.213:9092"
	defaultNatsUrl        = "nats://nats.jreed.svc.cluster.local:4222"
	defaultNsqConsumerUrl = "127.0.0.1:4161"
	defaultNsqProducerUrl = "127.0.0.1:4150"
	testData              = "Test Suite!"
	subscriberTestData    = "Subscriber-Specific Message!"
)

type MQSuiteConfig struct {
	urls              string
	subTopics         []string
	pubTopics         []string
	implementation    string
	deliverySemantics string

	consumerUrl string
	producerUrl string
}

func NewMQSuiteConfig(
	urls string,
	subTopics []string,
	pubTopics []string,
	implementation string,
	deliverySemantics string,
	consumerUrl string,
	producerUrl string) MQSuiteConfig {

	return MQSuiteConfig{
		urls,
		subTopics,
		pubTopics,
		implementation,
		deliverySemantics,
		consumerUrl,
		producerUrl,
	}
}

func MQSuite(config MQSuiteConfig) {

	fmt.Println("		       ┌───────────────────────────┐")
	fmt.Println("		       │     MQ Test Suite         │")
	fmt.Println("		       └───────────────────────────┘")

	// Create a new logger
	var logger *zap.Logger
	var err error
	if logger, err = logging.NewLogger(logging.Config{Type: "development"}); err != nil {
		fmt.Println("Could not set up logger for MQ CLI Test Suite. Exiting ...")
		return
	}

	// Create a new message queue
	fmt.Println("Creating a new", config.implementation, "Message Queue ...")
	var queue mq.MessageQueue
	var timeout time.Duration

	// Choose a different queue based on implementation
	switch config.implementation {
	case "kafka":
		// Create a unique namespace for test to avoid conflicting with other CLI suites
		uuidNamespace := uuid.New()
		mq.SetNamespace("cli-test-" + uuidNamespace.String())

		var kafkaUrls []string
		if config.urls == "" {
			config.urls = defaultKafkaUrl
			kafkaUrls = append(kafkaUrls, config.urls)
		}
		if queue, err = kafka.NewKafkaMessageQueue(logger, kafkaUrls); err != nil {
			fmt.Println("Unable to create a new Kafka Message Queue. Returned error:")
			errEncountered(err.Error())
			return
		}
		// Timeout for the handler to receive messages
		// If all messages are not received within this window, the test fails
		timeout = 11 * time.Second

	case "local":
		if queue, err = mock.NewLocalMessageQueue(logger, "development"); err != nil {
			fmt.Println("Unable to create a new Local Message Queue. Returned error:")
			errEncountered(err.Error())
			return
		}
		// Mock implementation supports wildcards.
		// This line adds new topics containing wildcards to be subscribed to.
		config.subTopics = append(config.subTopics, "foo.*", "foo.>", "*", "*.*.*", ">", "*.>")

		// Timeout for the handler to receive messages
		// If all messages are not received within this window, the test fails
		timeout = 5 * time.Second

	case "nats":
		// Create a unique namespace for test to avoid conflicting with other CLI suites
		uuidNamespace := uuid.New()
		mq.SetNamespace("cli-test-" + uuidNamespace.String())

		if config.urls == "" {
			config.urls = defaultNatsUrl
		}
		if queue, err = nats.NewNatsMessageQueue(logger, config.urls); err != nil {
			fmt.Println("Unable to create a new NATS Message Queue. Returned error:")
			errEncountered(err.Error())
			return
		}
		// NATS implementation supports wildcards.
		// This line adds new topics containing wildcards to be subscribed to.
		config.subTopics = append(config.subTopics, "foo.*", "foo.>", "*", "*.*.*", ">", "*.>")

		// Timeout for the handler to receive messages
		// If all messages are not received within this window, the test fails
		timeout = 5 * time.Second

	case "nsq":
		if config.consumerUrl == "" {
			config.consumerUrl = defaultNsqConsumerUrl
		}
		if config.producerUrl == "" {
			config.producerUrl = defaultNsqProducerUrl
		}
		if queue, err = nsq.NewNsqMessageQueue(logger, config.consumerUrl, config.producerUrl); err != nil {
			fmt.Println("Unable to create a new NATS Message Queue. Returned error:")
			errEncountered(err.Error())
			return
		}

		// Timeout for the handler to receive messages
		// If all messages are not received within this window, the test fails
		timeout = 5 * time.Second

	// If an invalid implementation is specified
	default:
		errEncountered(mq.ErrInvalidImplementation.Error())
		fmt.Println("Valid options include: kafka, local, nats.")
		return
	}

	fmt.Println("Successfully created a new", config.implementation, "Message Queue!")

	subs, msgExpectedCount := subscribeAll(config, queue)

	// Publish to the provided topics
	fmt.Println("Publishing to the provided topics ...")
	fmt.Println("Topics containing 'bar' will receive a subscriber-specific message.")
	for _, topic := range config.pubTopics {
		if strings.Contains(topic, "bar") {
			err = publish(queue, topic, subscriberTestData)
		} else {
			err = publish(queue, topic, testData)
			if err != nil {
				fmt.Println("Unable to publish to the provided topic", topic, "| Returned error:")
				errEncountered(err.Error())
				return
			}
		}
	}

	fmt.Println("Successfully published to all topics!")

	// Sleep while waiting for all messages to be received by their appropriate handlers
	fmt.Println("Checking if all messages were successfully received ...")
	time.Sleep(timeout)

	msgReceivedMutex.Lock()
	// If not all sent messages were received appropriately, fail the test
	if msgExpectedCount != msgReceivedCount {
		errEncountered(fmt.Sprintf(
			"%d messages received where %d expected",
			msgReceivedCount,
			msgExpectedCount,
		))
		msgReceivedMutex.Unlock()
		return
	}
	msgReceivedMutex.Unlock()

	fmt.Println("Successfully received all messages!")

	{
		// Test delayed publishing.
		msgReceivedMutex.Lock()
		msgReceivedCount = 0
		msgReceivedMutex.Unlock()

		var expectedErr error
		switch config.implementation {
		case "kafka":
		case "nats":
			msgExpectedCount = 0
			expectedErr = mq.ErrDelayedPublishUnsupported
		case "mock":
		case "nsq":
			expectedErr = nil
		}

		fmt.Println("Publishing to the provided topics on delay...")
		for _, topic := range config.pubTopics {
			err = publishDelay(queue, topic, "Delayed publish.")
			if err != expectedErr {
				fmt.Println("Incorrect error when doing delayed publish...")
				errEncountered("Delayed publish failure.")
				return
			}
		}

		// Sleep for half the delay time.
		fmt.Println("Waiting for message queue to be cleared ...")
		time.Sleep(5 * time.Second)

		// Confirm nothing received yet.
		msgReceivedMutex.Lock()
		if msgReceivedCount != 0 {
			errEncountered("Delayed publish receiving too soon.")
			msgReceivedMutex.Unlock()
			return
		}
		msgReceivedMutex.Unlock()

		// Sleep for delay time.
		fmt.Println("Waiting for message queue to be cleared ...")
		time.Sleep(10 * time.Second)

		// Confirm received.
		msgReceivedMutex.Lock()
		if msgReceivedCount != msgExpectedCount {
			errEncountered(fmt.Sprintf(
				"%d delayed messages received where %d expected",
				msgReceivedCount,
				msgExpectedCount,
			))
			msgReceivedMutex.Unlock()
			return
		}
		msgReceivedMutex.Unlock()

		fmt.Println("Successfully received all delayed messages!")
	}

	// Unsubscribe from the provided topics
	fmt.Println("Unsubscribing from topics ...")
	if unsubscribeAll(subs) != nil {
		return
	}
	fmt.Println("Successfully unsubscribed from all topics!")

	// Publish to the provided topics once more to verify no more messages are being received
	// Reset the message received count to accomplish this
	msgReceivedMutex.Lock()
	msgReceivedCount = 0
	msgReceivedMutex.Unlock()

	fmt.Println("Publishing to the provided topics post-unsubscription ...")
	for _, topic := range config.pubTopics {
		err = publish(queue, topic, "Published after unsubscription - you shouldn't be receiving this!")
		if err != nil {
			fmt.Println("Unable to publish to the provided topics. Returned error:")
			errEncountered(err.Error())
			return
		}
	}

	msgReceivedMutex.Lock()
	if msgReceivedCount != 0 {
		fmt.Println("Unable to publish to the provided topics. Returned error:")
		errEncountered("Message handlers are still receiving messages after unsubscribing from all topics.")
		msgReceivedMutex.Unlock()
		return
	}
	msgReceivedMutex.Unlock()

	// Resubscribe, receive, and unsubscribe to clear message queue of the messages we just published.
	subs, msgExpectedCount = subscribeAll(config, queue)

	// Sleep while waiting for all messages to be received by their appropriate handlers
	fmt.Println("Waiting for message queue to be cleared ...")
	time.Sleep(timeout)

	if unsubscribeAll(subs) != nil {
		return
	}

	fmt.Println("Successfully published to all topics!")
	fmt.Println("		       ┌────────────────────────────┐")
	fmt.Println("		       │    MQ Test Result: PASS    │")
	fmt.Println("		       └────────────────────────────┘")
}

func errEncountered(err string) {
	fmt.Println("- - - - -")
	fmt.Println(err)
	fmt.Println("- - - - -")
	fmt.Println("		       ┌────────────────────────────┐")
	fmt.Println("		       │    MQ Test Result: FAIL    │")
	fmt.Println("		       └────────────────────────────┘")
}

func subscribeAtMostOnce(
	queue mq.MessageQueue,
	topic string,
	topicNumber int,
) (mq.Subscription, error) {
	if sub, err := queue.Subscribe(topic, responseHandlerFactory(topicNumber), mq.WithAtMostOnceDelivery("test")); err != nil {
		return nil, err
	} else {
		fmt.Println("Successfully subscribed to topic", topic+"!")
		return sub, nil
	}
}

func subscribeAtLeastOnce(
	queue mq.MessageQueue,
	topic string,
	topicNumber int,
) (mq.Subscription, error) {
	if sub, err := queue.Subscribe(topic, responseHandlerFactory(topicNumber), mq.WithAtLeastOnceDelivery()); err != nil {
		return nil, err
	} else {
		fmt.Println("Successfully subscribed to topic", topic+"!")
		return sub, nil
	}
}

func publish(queue mq.MessageQueue, topic string, data string) (err error) {
	if err := queue.Publish(topic, mq.WithBytes([]byte(data))); err != nil {
		return err
	} else {
		fmt.Println("Successfully published topic", topic+"!")
	}
	return
}

func publishDelay(queue mq.MessageQueue, topic string, data string) (err error) {
	if err := queue.Publish(topic, mq.WithBytes([]byte(data)), mq.WithDelay(10*time.Second)); err != nil {
		return err
	} else {
		fmt.Println("Successfully published topic", topic+"!")
	}
	return
}

func responseHandlerFactory(subID int) mq.SubResponseHandler {
	return func(msg mq.Message, err error) mq.ConsumptionCode {
		if err != nil {
			if err == context.Canceled {
				return mq.ConsumeAckFinal
			} else {
				fmt.Println("handler received error:", err)
			}
		} else {
			fmt.Println("Msg received.",
				"SubID:", strconv.Itoa(subID),
				" || Topic:", msg.Topic(),
				" || Data:", string(msg.Data()))

			// Thread safety measure:
			// We don't want these handlers increasing the message received count simultaneously
			// Instead, we want this to happen sequentially.
			msgReceivedMutex.Lock()
			msgReceivedCount++
			msgReceivedMutex.Unlock()
		}
		return mq.ConsumeAck
	}
}

func subscribeAll(config MQSuiteConfig, queue mq.MessageQueue) (subs []mq.Subscription, msgExpectedCount int) {
	// Subscribe to the provided topics
	switch config.deliverySemantics {
	case "AtMostOnce":
		// Specify how many messages are expected to be received.
		switch config.implementation {
		case "kafka":
			msgExpectedCount = atMostOnceKafkaMsgExpectedCt
		case "local":
			msgExpectedCount = atMostOnceLocalMsgExpectedCt
		case "nats":
			msgExpectedCount = atMostOnceNatsMsgExpectedCt
		case "nsq":
			msgExpectedCount = atMostOnceNsqMsgExpectedCt
		default:
			msgExpectedCount = 1
		}

		fmt.Println("Subscribing to the provided topics with AtMostOnce ...")
		for i, topic := range config.subTopics {
			if sub, err := subscribeAtMostOnce(queue, topic, i); err != nil {
				fmt.Println("Unable to subscribe to the provided topic", topic, "| Returned error:")
				errEncountered(err.Error())
				return
			} else {
				subs = append(subs, sub)
			}
		}

	case "AtLeastOnce":
		// Specify how many messages are expected to be received.
		switch config.implementation {
		case "kafka":
			msgExpectedCount = atLeastOnceKafkaMsgExpectedCt
		case "local":
			msgExpectedCount = atLeastOnceLocalMsgExpectedCt
		case "nats":
			msgExpectedCount = atLeastOnceNatsMsgExpectedCt
		case "nsq":
			msgExpectedCount = atLeastOnceNsqMsgExpectedCt
		default:
			msgExpectedCount = 1
		}

		fmt.Println("Subscribing to the provided topics with AtLeastOnce ...")
		for i, topic := range config.subTopics {
			if sub, err := subscribeAtLeastOnce(queue, topic, i); err != nil {
				fmt.Println("Unable to subscribe to the provided topic", topic, "| Returned error:")
				errEncountered(err.Error())
				return nil, 0
			} else {
				// Add each subscription to a slice of subscriptions
				subs = append(subs, sub)
			}
		}

	// If invalid or unsupported delivery semantics are specified
	default:
		errEncountered(mq.ErrInvalidDeliverySemantics.Error())
		fmt.Println("Valid options include: AtMostOnce, AtLeastOnce.")
		return nil, 0
	}
	fmt.Println("Successfully subscribed to all topics!")
	return
}

func unsubscribeAll(subs []mq.Subscription) error {
	for _, sub := range subs {
		if err := sub.Unsubscribe(); err != nil {
			fmt.Println("Unable to unsubscribe from the provided topics. Returned error:")
			errEncountered(err.Error())
			return err
		}
	}

	return nil
}
