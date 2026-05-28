package nosql

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

// MockMessageQueue 模擬消息隊列
type MockMessageQueue struct {
	mock.Mock
}

func (m *MockMessageQueue) Publish(topic string, opts ...miface.PubOption) error {
	args := m.Called(topic, opts)
	return args.Error(0)
}

func (m *MockMessageQueue) Subscribe(ctx context.Context, topic string, handler miface.SubResponseHandler, opts ...miface.SubOption) (miface.Subscription, error) {
	args := m.Called(ctx, topic, handler, opts)
	return args.Get(0).(miface.Subscription), args.Error(1)
}

// MockSubscription 模擬訂閱
type MockSubscription struct {
	mock.Mock
}

func (m *MockSubscription) IsValid() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSubscription) Unsubscribe() error {
	args := m.Called()
	return args.Error(0)
}

// MockDocumentProvider 模擬文檔提供者
type MockDocumentProvider struct {
	mock.Mock
}

func (m *MockDocumentProvider) OpenDbDriver(collectionName string) (diface.ICollection, error) {
	args := m.Called(collectionName)
	return args.Get(0).(diface.ICollection), args.Error(1)
}

func (m *MockDocumentProvider) Shutdown() error {
	args := m.Called()
	return args.Error(0)
}

// MockCollection 模擬集合
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockCollection) Set(ctx context.Context, key key.Key, opts ...noptions.Option) (noptions.Version, error) {
	args := m.Called(ctx, key, opts)
	return args.Get(0).(noptions.Version), args.Error(1)
}

func (m *MockCollection) Get(ctx context.Context, key key.Key, opts ...noptions.Option) (noptions.Version, error) {
	args := m.Called(ctx, key, opts)
	return args.Get(0).(noptions.Version), args.Error(1)
}

func (m *MockCollection) Delete(ctx context.Context, key key.Key) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestWriteBackOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    WriteBackOptions
		wantErr bool
	}{
		{
			name: "valid disabled options",
			opts: WriteBackOptions{
				Enabled: false,
				Delay:   time.Second,
				MQ:      nil,
			},
			wantErr: false,
		},
		{
			name: "valid enabled options",
			opts: WriteBackOptions{
				Enabled: true,
				Delay:   time.Second,
				MQ:      &MockMessageQueue{},
			},
			wantErr: false,
		},
		{
			name: "invalid: enabled but no MQ",
			opts: WriteBackOptions{
				Enabled: true,
				Delay:   time.Second,
				MQ:      nil,
			},
			wantErr: true,
		},
		{
			name: "invalid: negative delay",
			opts: WriteBackOptions{
				Enabled: false,
				Delay:   -time.Second,
				MQ:      nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteBackWorker_GetMetrics(t *testing.T) {
	mockMQ := &MockMessageQueue{}
	mockProvider := &MockDocumentProvider{}
	logger := zap.NewNop()

	worker := NewWriteBackWorker(mockMQ, mockProvider, logger)
	metrics := worker.GetMetrics()

	assert.NotNil(t, metrics)
	assert.GreaterOrEqual(t, metrics.ProcessedCount, int64(0))
	assert.GreaterOrEqual(t, metrics.FailedCount, int64(0))
}

func TestWriteBackPayload_JSON(t *testing.T) {
	payload := WriteBackPayload{
		CollectionName: "test_collection",
		Key:           "test_key",
		Data: map[string]any{
			"name": "John",
			"age":  30,
		},
		Version: noptions.Version(1),
	}

	// 測試序列化
	data, err := json.Marshal(payload)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// 測試反序列化
	var decoded WriteBackPayload
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, payload.CollectionName, decoded.CollectionName)
	assert.Equal(t, payload.Key, decoded.Key)
	assert.Equal(t, payload.Version, decoded.Version)
}
