package mockmq

import (
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/gstones/platform/services/common/mq"
)

func TestLocalSubscribeAndPublish(t *testing.T) {
	type testCase struct {
		subTopic    string // Topic to be subscribed to. Supports wildcards.
		pubTopic    string // Topic to be published to. Does not support wildcards.
		expectedErr error  // The error expected, if any.
	}

	testCases := []testCase{
		{
			subTopic:    "test.>",
			pubTopic:    "test",
			expectedErr: nil,
		},
		{
			subTopic:    "test.>",
			pubTopic:    "test",
			expectedErr: nil,
		},
		{
			subTopic:    "*.topic",
			pubTopic:    "test.topic",
			expectedErr: mq.ErrTopicValidationFail,
		},
		{
			subTopic:    "*",
			pubTopic:    "test",
			expectedErr: mq.ErrTopicValidationFail,
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
	logger := zaptest.NewLogger(t)

	for i, tc := range testCases {
		queue := NewMessageQueue(logger, "development")

		if sub, err := queue.Subscribe(
			tc.subTopic,
			messageHandlerFactory(),
		); err != nil {
			if err == tc.expectedErr {
				continue
			} else {
				t.Fatal("Error encountered in test case #", i+1, "in QueueSubscribe() call:", err)
			}
		} else {
			if sub.IsValid() == false {
				t.Fatal("Subscription generated in test case #", i+1, "is invalid.")
			} else {
				repeatCount := 1
				for r := 0; r < repeatCount; r++ {
					if err := queue.Publish(tc.pubTopic, []byte("test!")); err != nil {
						if err == mq.ErrEmptyTopic {
							continue
						} else {
							t.Fatal("Error encountered in test case", i+1, "in Publish() call:", err)
						}
					}
				}
				if err := sub.Unsubscribe(); err != nil {
					t.Fatal("Error encountered in test case", i+1, "in Unsubscribe() call:", err)
				}
			}
		}
	}
}

func TestQueueSubscribe(t *testing.T) {
	type testCase struct {
		subTopic    string // Topic to be subscribed to. Supports wildcards.
		queueName   string
		errExpected error
	}

	testCases := []testCase{
		{
			subTopic:    "test.>",
			queueName:   "queue",
			errExpected: nil,
		},
		{
			subTopic:    "*.topic",
			queueName:   "test",
			errExpected: mq.ErrTopicValidationFail,
		},
		{
			subTopic:    "*",
			queueName:   "test",
			errExpected: mq.ErrTopicValidationFail,
		},
		{
			subTopic:    ">.test.topic",
			queueName:   "test",
			errExpected: mq.ErrTopicValidationFail,
		},
		{
			subTopic:    "",
			errExpected: mq.ErrEmptyTopic,
		},
	}

	logger := zaptest.NewLogger(t)
	queue := NewMessageQueue(logger, "development")
	for i, tc := range testCases {
		if sub, err := queue.QueueSubscribe(
			tc.subTopic,
			tc.queueName,
			messageHandlerFactory(),
		); err != nil {
			if err == tc.errExpected {
				continue
			} else {
				t.Fatal("Error encountered in test case #", i+1, "in QueueSubscribe() call:", err)
			}
		} else {
			if sub.Queue() != tc.queueName {
				t.Fatal("Subscription queue field does not match provided queue name in test case #", i+1)
			} else if sub.IsValid() == false {
				t.Fatal("Generated subscription in test case #", i+1, "is not valid.")
			}
		}
	}
}
