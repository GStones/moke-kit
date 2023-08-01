package mockmq

import (
	"testing"

	"github.com/gstones/platform/services/common/mq"
)

func TestSubscription(t *testing.T) {
	type testCase struct {
		queue       string
		message     Message
		handler     MessageHandler
		errExpected error
	}
	testCases := []testCase{
		{
			"Test",
			Message{"testTopic", []byte("testMessage!")},
			messageHandlerFactory(),
			nil,
		},
		{
			"Queue",
			Message{"queueTopic", []byte("queueMessage!")},
			messageHandlerFactory(),
			nil,
		},
		{
			"Queue",
			Message{"queueTopic", []byte("queueMessage!")},
			nil,
			mq.ErrInvalidSubscription,
		},
	}
	handler := messageHandlerFactory()
	subGroup := NewPermSubGroup()

	for i, tc := range testCases {
		if sub := subGroup.NewSubscription(tc.queue, handler); sub.IsValid() == false {
			t.Fatal("Generated subscription is not valid in test case #", i+1)
		} else {
			if sub.Queue() != tc.queue {
				t.Fatal("Subscription queue name is not equal to the test case queue's name in test case #", i+1)
			}
			if err := sub.SendMessage(tc.message); err != nil {
				t.Fatal("Could not send message in test case #", i+1, ":", err)
			} else {

				if err := sub.Unsubscribe(); tc.errExpected == mq.ErrInvalidSubscription && err == tc.errExpected {
					continue
				} else if err != nil {
					t.Fatal("Error encountered unsubscribing from subscription queue",
						sub.Queue(),
						"in test case #", i+1, ":", err)
				}
			}
		}
	}
}
func messageHandlerFactory() MessageHandler {
	return func(msg Message) {
	}
}
