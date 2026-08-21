package nosql

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nosql/diface"
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

func TestWriteBackWorker_ApplyNewerTargetVersion(t *testing.T) {
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

	// Newer target arrives first (rapid SaveAsync / reordered MQ).
	assert.NoError(t, mq.Publish(WriteBackTopic, miface.WithJSON(&WriteBackPayload{
		CollectionName: "wb-rebase",
		Key:            k.String(),
		Data:           json.RawMessage(`{"message":"latest"}`),
		Version:        5,
	})))
	assert.NoError(t, mq.Publish(WriteBackTopic, miface.WithJSON(&WriteBackPayload{
		CollectionName: "wb-rebase",
		Key:            k.String(),
		Data:           json.RawMessage(`{"message":"older"}`),
		Version:        2,
	})))

	assert.Eventually(t, func() bool {
		return worker.GetMetrics().ProcessedCount >= 1 && worker.GetMetrics().FailedCount >= 1
	}, time.Second, 10*time.Millisecond)

	var got docPayload
	ver, err := coll.Get(context.Background(), k, noptions.WithDestination(&got))
	assert.NoError(t, err)
	assert.Equal(t, "latest", got.Message)
	// Fast-forward must land exactly on the optimistic target version.
	assert.Equal(t, noptions.Version(5), ver)
}

type jumpOnSecondGetCollection struct {
	diface.ICollection
	gets int
}

func (c *jumpOnSecondGetCollection) Get(
	ctx context.Context,
	k key.Key,
	opts ...noptions.Option,
) (noptions.Version, error) {
	c.gets++
	if c.gets == 2 {
		// Between fast-forward iterations, an external writer advances the chain.
		ver, err := c.ICollection.Get(ctx, k, opts...)
		if err != nil {
			return ver, err
		}
		_, err = c.Set(
			ctx,
			k,
			noptions.WithSource(&docPayload{Message: "external"}),
			noptions.WithVersion(ver),
		)
		if err != nil {
			return 0, err
		}
	}
	return c.ICollection.Get(ctx, k, opts...)
}

func TestApplyWriteBackSnapshot_AbortsOnExternalVersionJump(t *testing.T) {
	logger := zap.NewNop()
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("wb-jump")
	assert.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "jump-1")
	assert.NoError(t, err)
	_, err = coll.Set(context.Background(), k, noptions.WithSource(&docPayload{Message: "v1"}))
	assert.NoError(t, err)

	wrapped := &jumpOnSecondGetCollection{ICollection: coll}
	err = applyWriteBackSnapshot(
		context.Background(),
		wrapped,
		k,
		map[string]any{"message": "wb"},
		5,
		nil,
		0,
	)
	assert.ErrorIs(t, err, errWriteBackStale)

	var got docPayload
	ver, err := coll.Get(context.Background(), k, noptions.WithDestination(&got))
	assert.NoError(t, err)
	assert.Equal(t, "external", got.Message)
	assert.Equal(t, noptions.Version(3), ver)
}

func TestWriteBackWorker_OutOfOrderTargetsKeepNewest(t *testing.T) {
	logger := zap.NewNop()
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("wb-ooo")
	assert.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "ooo-1")
	assert.NoError(t, err)

	_, err = coll.Set(context.Background(), k, noptions.WithSource(&docPayload{Message: "v1"}))
	assert.NoError(t, err)

	mq := &capturingMQ{}
	worker := NewWriteBackWorker(mq, provider, logger)
	assert.NoError(t, worker.Start())
	defer worker.Stop()

	assert.NoError(t, mq.Publish(WriteBackTopic, miface.WithJSON(&WriteBackPayload{
		CollectionName: "wb-ooo",
		Key:            k.String(),
		Data:           json.RawMessage(`{"message":"latest"}`),
		Version:        5,
	})))
	assert.NoError(t, mq.Publish(WriteBackTopic, miface.WithJSON(&WriteBackPayload{
		CollectionName: "wb-ooo",
		Key:            k.String(),
		Data:           json.RawMessage(`{"message":"older"}`),
		Version:        4,
	})))

	assert.Eventually(t, func() bool {
		return worker.GetMetrics().ProcessedCount >= 1 && worker.GetMetrics().FailedCount >= 1
	}, time.Second, 10*time.Millisecond)

	var got docPayload
	ver, err := coll.Get(context.Background(), k, noptions.WithDestination(&got))
	assert.NoError(t, err)
	assert.Equal(t, "latest", got.Message)
	assert.Equal(t, noptions.Version(5), ver)
}

func TestWriteBackWorker_RejectEpochMismatch(t *testing.T) {
	logger := zap.NewNop()
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("wb-epoch")
	assert.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "epoch-1")
	assert.NoError(t, err)

	// Recreated document at version 1 (old delayed write-back may carry a higher target).
	_, err = coll.Set(context.Background(), k, noptions.WithSource(&docPayload{Message: "recreated"}))
	assert.NoError(t, err)

	cache := newMemoryHashCache()
	assert.NoError(t, cache.SetCache(context.Background(), k, map[string]any{
		cacheFieldVersion: noptions.Version(1),
		cacheFieldEpoch:   int64(2),
		cacheFieldData:    `{"message":"recreated"}`,
	}, time.Minute))

	mq := &capturingMQ{}
	worker := NewWriteBackWorker(mq, provider, logger).WithCache(cache)
	assert.NoError(t, worker.Start())
	defer worker.Stop()

	assert.NoError(t, mq.Publish(WriteBackTopic, miface.WithJSON(&WriteBackPayload{
		CollectionName: "wb-epoch",
		Key:            k.String(),
		Data:           json.RawMessage(`{"message":"stale-generation"}`),
		Version:        9,
		Epoch:          1,
	})))

	assert.Eventually(t, func() bool {
		return worker.GetMetrics().FailedCount >= 1
	}, time.Second, 10*time.Millisecond)

	var got docPayload
	_, err = coll.Get(context.Background(), k, noptions.WithDestination(&got))
	assert.NoError(t, err)
	assert.Equal(t, "recreated", got.Message)
	assert.Equal(t, int64(0), worker.GetMetrics().ProcessedCount)
}

func TestWriteBackWorker_RejectStaleTargetVersion(t *testing.T) {
	logger := zap.NewNop()
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("wb-stale")
	assert.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "stale-1")
	assert.NoError(t, err)

	_, err = coll.Set(context.Background(), k, noptions.WithSource(&docPayload{Message: "seed"}))
	assert.NoError(t, err)
	_, err = coll.Set(
		context.Background(),
		k,
		noptions.WithSource(&docPayload{Message: "fresh"}),
		noptions.WithVersion(1),
	)
	assert.NoError(t, err)

	mq := &capturingMQ{}
	worker := NewWriteBackWorker(mq, provider, logger)
	assert.NoError(t, worker.Start())
	defer worker.Stop()

	assert.NoError(t, mq.Publish(WriteBackTopic, miface.WithJSON(&WriteBackPayload{
		CollectionName: "wb-stale",
		Key:            k.String(),
		Data:           json.RawMessage(`{"message":"stale"}`),
		Version:        2, // not newer than current DB version 2
	})))

	assert.Eventually(t, func() bool {
		return worker.GetMetrics().FailedCount >= 1
	}, time.Second, 10*time.Millisecond)

	var got docPayload
	_, err = coll.Get(context.Background(), k, noptions.WithDestination(&got))
	assert.NoError(t, err)
	assert.Equal(t, "fresh", got.Message)
	assert.Equal(t, int64(0), worker.GetMetrics().ProcessedCount)
}

func TestWriteBackPayload_JSON(t *testing.T) {
	payload := WriteBackPayload{
		CollectionName: "test_collection",
		Key:            "test_key",
		Data:           json.RawMessage(`{"name":"John","age":30}`),
		Version:        noptions.Version(1),
		Epoch:          3,
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
	assert.Equal(t, payload.Epoch, decoded.Epoch)
}
