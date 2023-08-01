package internal

import (
	"testing"
	"time"

	"go.uber.org/zap/zaptest"

	"github.com/gstones/platform/services/common/mq"
)

const (
	receiveTimeout = 4 * time.Second
)

func TestLocalSubscribeAndPublish(t *testing.T) {
	type testCase struct {
		subTopic          string // Topic to be subscribed to. Supports wildcards.
		pubTopic          string // Topic to be published to. Does not support wildcards.
		groupId           string
		data              string               // Data to be sent with the publish.
		deliverySemantics mq.DeliverySemantics // AtMostOnce, AtLeastOnce, etc.
		expectedErr       error                // The error expected, if any.
	}

	testCases := []testCase{
		{
			subTopic:          "test.>",
			pubTopic:          "test.case.topic",
			groupId:           "test",
			data:              "test case 1 publish data",
			deliverySemantics: mq.AtMostOnce,
			expectedErr:       nil,
		},
		{
			subTopic:          "test.topic.>",
			pubTopic:          "test.topic.case",
			groupId:           "case",
			data:              "test case 2 publish data",
			deliverySemantics: mq.AtLeastOnce,
			expectedErr:       nil,
		},
		{
			subTopic:          "*.topic",
			pubTopic:          "test.topic",
			data:              "test case 3 publish data",
			deliverySemantics: mq.AtMostOnce,
			expectedErr:       nil,
		},
		{
			subTopic:          "*",
			pubTopic:          "test",
			data:              "test case 4 publish data",
			deliverySemantics: mq.AtLeastOnce,
			expectedErr:       nil,
		},
		{
			subTopic:    ">.test.topic",
			pubTopic:    "test",
			expectedErr: mq.ErrTopicValidationFail,
		},
		{
			subTopic:    "",
			expectedErr: mq.ErrEmptyTopic,
		},
	}

	var subOpt mq.SubOption
	logger := zaptest.NewLogger(t)

	for i, tc := range testCases {
		// Define delivery semantics options
		atMostOnceOpt := mq.WithAtMostOnceDelivery(mq.GroupId(tc.groupId))
		atLeastOnceOpt := mq.WithAtLeastOnceDelivery()

		// Create a new message queue
		if queue, err := NewMessageQueue(logger, "development"); err != nil {
			t.Fatal("Error encountered creating a new Message Queue", err)
		} else {
			switch tc.deliverySemantics {
			case mq.AtLeastOnce:
				subOpt = atLeastOnceOpt
			case mq.AtMostOnce:
				subOpt = atMostOnceOpt
			}

			// Set up handler response channel and handler error channel
			responseChannel := make(chan mq.Message)
			errorChannel := make(chan error)
			// Subscribe to the given topic
			if sub, err := queue.Subscribe(
				tc.subTopic,
				responseHandlerFactory(responseChannel, errorChannel),
				subOpt,
			); err != nil {
				if err == tc.expectedErr {
					continue
				} else {
					t.Fatal("Unexpected error encountered in test case #", i+1, ":", err)
				}
			} else if sub.IsValid() == false {
				t.Fatal("Subscription generated in test case #", i+1, "is invalid.")
			} else if err == nil && err != tc.expectedErr {
				t.Fatal("Error was not encountered where one was expected in test case #", i+1)
			} else {
				var recvCt int

				// Publish to the given topic
				go func() {
					if err := queue.Publish(tc.pubTopic, mq.WithBytes([]byte(tc.data))); err != nil {
						t.Fatal("Unexpected error encountered in test case #", i+1, "in Publish():", err)
					}
				}()

				// Start receiving on the channel
			recvLoop:
				for {
					select {
					case err := <-errorChannel:
						t.Error("Test case", i+1, "encountered error:", err)
						break recvLoop
					case <-responseChannel:
						recvCt++
					case <-time.After(receiveTimeout):
						break recvLoop
					}
				}

				// Check if we've received the expected message - no more, no less
				if recvCt != 1 {
					t.Fatal("Message received count was not equal to the expected count in test case #", i+1)
				}

				// Unsubscribe from the current subscription
				if err := sub.Unsubscribe(); err != nil {
					t.Fatal("Error encountered in test case", i+1, "in Unsubscribe() call:", err)
				}
			}
		}
	}
}

func responseHandlerFactory(responseChannel chan mq.Message, errorChannel chan error) mq.SubResponseHandler {
	return func(msg mq.Message, err error) mq.ConsumptionCode {
		if err != nil {
			errorChannel <- err
		} else {
			responseChannel <- msg
		}
		return mq.ConsumeAck
	}
}
