package internal

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	type testCase struct {
		id    string
		topic string
		data  []byte
		vPtr  interface{}
	}
	testCases := []testCase{
		{
			id:    "testId",
			topic: "testTopic",
			data:  []byte("Test!"),
			vPtr:  nil,
		},
	}
	for _, tc := range testCases {
		message := NewMessage(tc.id, tc.topic, tc.data, tc.vPtr)
		if message.ID() != tc.id {
			t.Fatal("Test Message's ID field does not match the provided test case.")
		}
		if message.Topic() != tc.topic {
			t.Fatal("Test Message's topic field does not match the provided test case.")
		}
		if string(message.Data()) != string(tc.data) {
			t.Fatal("Test Message's data field does not match the provided test case.")
		}
		if message.VPtr() != tc.vPtr {
			t.Fatal("Test Message's vPtr field does not match the provided test case.")
		}
	}
}
