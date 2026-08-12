package nosql

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/mock"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

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
				MQ:      &capturingMQ{},
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
	worker := NewWriteBackWorker(
		&capturingMQ{},
		mock.NewMockDriverProvider(zap.NewNop()),
		zap.NewNop(),
	)
	metrics := worker.GetMetrics()
	assert.GreaterOrEqual(t, metrics.ProcessedCount, int64(0))
	assert.GreaterOrEqual(t, metrics.FailedCount, int64(0))
}

func TestWriteBackWorker_InvalidKeyNoPanic(t *testing.T) {
	mq := &capturingMQ{}
	provider := mock.NewMockDriverProvider(zap.NewNop())
	worker := NewWriteBackWorker(mq, provider, zap.NewNop())
	assert.NoError(t, worker.Start())
	defer worker.Stop()

	assert.NoError(t, mq.Publish(WriteBackTopic, miface.WithJSON(&WriteBackPayload{
		CollectionName: "c",
		Key:            "not-a-valid-key",
		Data:           json.RawMessage(`{}`),
		Version:        1,
	})))

	// Invalid keys are nacked as persistent failures and must not crash the worker.
	assert.Equal(t, int64(0), worker.GetMetrics().ProcessedCount)
}

type failSubscribeMQ struct{}

func (m *failSubscribeMQ) Publish(topic string, opts ...miface.PubOption) error {
	return nil
}

func (m *failSubscribeMQ) Subscribe(
	ctx context.Context,
	topic string,
	handler miface.SubResponseHandler,
	opts ...miface.SubOption,
) (miface.Subscription, error) {
	return nil, assert.AnError
}

func TestWriteBackWorker_StartSubscribeError(t *testing.T) {
	worker := NewWriteBackWorker(
		&failSubscribeMQ{},
		mock.NewMockDriverProvider(zap.NewNop()),
		zap.NewNop(),
	)
	assert.Error(t, worker.Start())
}

func TestWriteBackWorker_RebaseOnVersionGap(t *testing.T) {
	logger := zap.NewNop()
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("wb-rebase")
	assert.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "rebase-1")
	assert.NoError(t, err)

	_, err = coll.Set(context.Background(), k, noptions.WithSource(&docPayload{Message: "v1"}))
	assert.NoError(t, err)

	mq := &capturingMQ{}
	worker := NewWriteBackWorker(mq, provider, logger)
	assert.NoError(t, worker.Start())
	defer worker.Stop()

	assert.NoError(t, mq.Publish(WriteBackTopic, miface.WithJSON(&WriteBackPayload{
		CollectionName: "wb-rebase",
		Key:            k.String(),
		Data:           json.RawMessage(`{"message":"rebased"}`),
		Version:        5, // gap: earlier write-backs never arrived
	})))

	assert.Eventually(t, func() bool {
		return worker.GetMetrics().ProcessedCount >= 1
	}, time.Second, 10*time.Millisecond)

	var got docPayload
	ver, err := coll.Get(context.Background(), k, noptions.WithDestination(&got))
	assert.NoError(t, err)
	assert.Equal(t, "rebased", got.Message)
	assert.Equal(t, noptions.Version(2), ver)
}

func TestWriteBackPayload_JSON(t *testing.T) {
	payload := WriteBackPayload{
		CollectionName: "test_collection",
		Key:            "test_key",
		Data:           json.RawMessage(`{"name":"John","age":30}`),
		Version:        noptions.Version(1),
	}

	data, err := json.Marshal(payload)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded WriteBackPayload
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, payload.CollectionName, decoded.CollectionName)
	assert.Equal(t, payload.Key, decoded.Key)
	assert.Equal(t, payload.Version, decoded.Version)
}
