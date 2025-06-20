package mq

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/test/utils"
)

// MockSubscription is a mock implementation of Subscription
type MockSubscription struct {
	mock.Mock
	valid bool
}

func (m *MockSubscription) IsValid() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSubscription) Unsubscribe() error {
	args := m.Called()
	return args.Error(0)
}

// MockMessage is a mock implementation of Message
type MockMessage struct {
	mock.Mock
	id    string
	topic string
	data  []byte
	vPtr  interface{}
}

func (m *MockMessage) ID() string {
	if m.id != "" {
		return m.id
	}
	args := m.Called()
	return args.String(0)
}

func (m *MockMessage) Topic() string {
	if m.topic != "" {
		return m.topic
	}
	args := m.Called()
	return args.String(0)
}

func (m *MockMessage) Data() []byte {
	if m.data != nil {
		return m.data
	}
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *MockMessage) VPtr() interface{} {
	if m.vPtr != nil {
		return m.vPtr
	}
	args := m.Called()
	return args.Get(0)
}

// MockMessageQueue is a mock implementation of MessageQueue
type MockMessageQueue struct {
	mock.Mock
	subscriptions map[string][]miface.SubResponseHandler
	mu            sync.RWMutex
}

func NewMockMessageQueue() *MockMessageQueue {
	return &MockMessageQueue{
		subscriptions: make(map[string][]miface.SubResponseHandler),
	}
}

func (m *MockMessageQueue) Subscribe(ctx context.Context, topic string, handler miface.SubResponseHandler, opts ...miface.SubOption) (miface.Subscription, error) {
	args := m.Called(ctx, topic, handler, opts)
	
	// Store handler for publish simulation
	m.mu.Lock()
	m.subscriptions[topic] = append(m.subscriptions[topic], handler)
	m.mu.Unlock()
	
	return args.Get(0).(miface.Subscription), args.Error(1)
}

func (m *MockMessageQueue) Publish(topic string, opts ...miface.PubOption) error {
	args := m.Called(topic, opts)
	
	// Simulate message delivery to subscribers
	m.mu.RLock()
	handlers := m.subscriptions[topic]
	m.mu.RUnlock()
	
	if len(handlers) > 0 {
		options, _ := miface.NewPubOptions(opts...)
		msg := &MockMessage{
			id:    "test-message-id",
			topic: topic,
			data:  options.Data,
		}
		
		for _, handler := range handlers {
			go handler(msg, nil)
		}
	}
	
	return args.Error(0)
}

// TestMessageInterfaces tests the message interfaces
func TestMessageInterfaces(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("PubOptions", func(t *testing.T) {
		testData := []byte("test data")
		delay := 5 * time.Second

		// Test WithBytes option
		opts, err := miface.NewPubOptions(miface.WithBytes(testData))
		helper.RequireNoError(err)
		helper.AssertEqual(testData, opts.Data)

		// Test WithDelay option
		opts, err = miface.NewPubOptions(miface.WithDelay(delay))
		helper.RequireNoError(err)
		helper.AssertEqual(delay, opts.Delay)

		// Test combined options
		opts, err = miface.NewPubOptions(
			miface.WithBytes(testData),
			miface.WithDelay(delay),
		)
		helper.RequireNoError(err)
		helper.AssertEqual(testData, opts.Data)
		helper.AssertEqual(delay, opts.Delay)
	})

	t.Run("SubOptions", func(t *testing.T) {
		opts, err := miface.NewSubOptions()
		helper.RequireNoError(err)
		helper.AssertNotNil(opts)

		// Test with AtLeastOnce delivery
		opts, err = miface.NewSubOptions(miface.WithAtLeastOnceDelivery())
		helper.RequireNoError(err)
		helper.AssertEqual(common.AtLeastOnce, opts.DeliverySemantics)

		// Test with AtMostOnce delivery
		groupId := common.GroupId("test-group")
		opts, err = miface.NewSubOptions(miface.WithAtMostOnceDelivery(groupId))
		helper.RequireNoError(err)
		helper.AssertEqual(common.AtMostOnce, opts.DeliverySemantics)
		helper.AssertEqual(string(groupId), opts.GroupId)
	})
}

// TestMockMessageQueue tests the mock message queue implementation
func TestMockMessageQueue(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	mockMQ := NewMockMessageQueue()
	mockSub := &MockSubscription{}

	t.Run("Subscribe", func(t *testing.T) {
		topic := "test-topic"
		handler := func(msg miface.Message, err error) common.ConsumptionCode {
			return common.ConsumeAck
		}

		mockSub.On("IsValid").Return(true)
		mockMQ.On("Subscribe", mock.Anything, topic, mock.Anything, mock.Anything).Return(mockSub, nil)

		sub, err := mockMQ.Subscribe(helper.Context(), topic, handler)
		helper.RequireNoError(err)
		helper.AssertNotNil(sub)
		helper.AssertTrue(sub.IsValid())

		mockMQ.AssertExpectations(t)
		mockSub.AssertExpectations(t)
	})

	t.Run("Publish", func(t *testing.T) {
		topic := "test-topic"
		testData := []byte("test message")

		mockMQ.On("Publish", topic, mock.Anything).Return(nil)

		err := mockMQ.Publish(topic, miface.WithBytes(testData))
		helper.AssertNoError(err)

		mockMQ.AssertExpectations(t)
	})

	t.Run("PublishSubscribeIntegration", func(t *testing.T) {
		topic := "integration-test"
		testData := []byte("integration message")
		received := make(chan []byte, 1)

		// Setup handler
		handler := func(msg miface.Message, err error) common.ConsumptionCode {
			helper.AssertNoError(err)
			received <- msg.Data()
			return common.ConsumeAck
		}

		// Mock subscription
		mockSub.On("IsValid").Return(true)
		mockMQ.On("Subscribe", mock.Anything, topic, mock.Anything, mock.Anything).Return(mockSub, nil)
		mockMQ.On("Publish", topic, mock.Anything).Return(nil)

		// Subscribe
		_, err := mockMQ.Subscribe(helper.Context(), topic, handler)
		helper.RequireNoError(err)

		// Publish
		err = mockMQ.Publish(topic, miface.WithBytes(testData))
		helper.RequireNoError(err)

		// Wait for message
		select {
		case receivedData := <-received:
			helper.AssertEqual(testData, receivedData)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for message")
		}

		mockMQ.AssertExpectations(t)
		mockSub.AssertExpectations(t)
	})
}

// TestConsumptionCodes tests consumption code handling
func TestConsumptionCodes(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("ConsumptionCodes", func(t *testing.T) {
		codes := []common.ConsumptionCode{
			common.ConsumeAck,
			common.ConsumeAckFinal,
			common.ConsumeNackTransientFailure,
			common.ConsumeNackPersistentFailure,
		}

		for _, code := range codes {
			helper.AssertNotNil(code, "Consumption code should not be nil")
		}
	})
}

// TestTopicHeaders tests topic header functionality
func TestTopicHeaders(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("LocalHeader", func(t *testing.T) {
		topic := "test-topic"
		fullTopic := common.LocalHeader.CreateTopic(topic)
		helper.AssertTrue(len(fullTopic) > len(topic), "Full topic should be longer than base topic")
	})

	t.Run("NatsHeader", func(t *testing.T) {
		topic := "test-topic"
		fullTopic := common.NatsHeader.CreateTopic(topic)
		helper.AssertTrue(len(fullTopic) > len(topic), "Full topic should be longer than base topic")
	})
}

// BenchmarkMockMessageQueue benchmarks the mock message queue
func BenchmarkMockMessageQueue(b *testing.B) {
	mockMQ := NewMockMessageQueue()
	mockSub := &MockSubscription{}
	topic := "benchmark-topic"
	
	// Setup mocks
	mockSub.On("IsValid").Return(true)
	mockMQ.On("Subscribe", mock.Anything, topic, mock.Anything, mock.Anything).Return(mockSub, nil)
	mockMQ.On("Publish", topic, mock.Anything).Return(nil)

	// Setup subscriber
	received := make(chan []byte, b.N)
	handler := func(msg miface.Message, err error) common.ConsumptionCode {
		received <- msg.Data()
		return common.ConsumeAck
	}
	
	_, err := mockMQ.Subscribe(context.Background(), topic, handler)
	require.NoError(b, err)

	testData := []byte("benchmark message")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := mockMQ.Publish(topic, miface.WithBytes(testData))
			if err != nil {
				b.Error(err)
			}
		}
	})

	// Wait for all messages
	for i := 0; i < b.N; i++ {
		select {
		case <-received:
		case <-time.After(10 * time.Second):
			b.Fatal("Timeout waiting for messages")
		}
	}
}
