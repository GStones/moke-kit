package nosql

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/mock"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

type docPayload struct {
	Message string `json:"message" bson:"message"`
}

type testDoc struct {
	DocumentBase
	ID   string
	Data *docPayload
}

type memoryHashCache struct {
	mu   sync.Mutex
	data map[string]map[string]any
}

func newMemoryHashCache() *memoryHashCache {
	return &memoryHashCache{data: make(map[string]map[string]any)}
}

func (c *memoryHashCache) GetCache(ctx context.Context, k key.Key, fields ...string) map[string]any {
	c.mu.Lock()
	defer c.mu.Unlock()
	src, ok := c.data[k.String()]
	if !ok {
		return nil
	}
	if len(fields) == 0 {
		out := make(map[string]any, len(src))
		for key, val := range src {
			out[key] = val
		}
		return out
	}
	out := make(map[string]any, len(fields))
	for _, field := range fields {
		if val, ok := src[field]; ok {
			out[field] = val
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func (c *memoryHashCache) SetCache(ctx context.Context, k key.Key, data map[string]any, expire time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	dst := make(map[string]any, len(data))
	for key, val := range data {
		dst[key] = val
	}
	c.data[k.String()] = dst
	return nil
}

func (c *memoryHashCache) DeleteCache(ctx context.Context, k key.Key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, k.String())
}

type capturingMQ struct {
	mu       sync.Mutex
	messages [][]byte
	handler  miface.SubResponseHandler
}

func (m *capturingMQ) Publish(topic string, opts ...miface.PubOption) error {
	options, err := miface.NewPubOptions(opts...)
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.messages = append(m.messages, append([]byte(nil), options.Data...))
	handler := m.handler
	m.mu.Unlock()
	if handler != nil {
		_ = handler(&staticMessage{data: options.Data}, nil)
	}
	return nil
}

func (m *capturingMQ) Subscribe(ctx context.Context, topic string, handler miface.SubResponseHandler, opts ...miface.SubOption) (miface.Subscription, error) {
	m.mu.Lock()
	m.handler = handler
	m.mu.Unlock()
	return &noopSub{}, nil
}

type staticMessage struct {
	data []byte
}

func (m *staticMessage) Topic() string { return WriteBackTopic }
func (m *staticMessage) Data() []byte  { return m.data }
func (m *staticMessage) ID() string    { return "1" }
func (m *staticMessage) VPtr() any     { return nil }

type noopSub struct{}

func (s *noopSub) IsValid() bool      { return true }
func (s *noopSub) Unsubscribe() error { return nil }

var (
	_ diface.ICache       = (*memoryHashCache)(nil)
	_ miface.MessageQueue = (*capturingMQ)(nil)
	_ miface.Message      = (*staticMessage)(nil)
	_ miface.Subscription = (*noopSub)(nil)
)

func TestDocumentBase_CRUD(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("doc-crud")
	require.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "10000")
	require.NoError(t, err)

	td := &testDoc{
		ID:   "10000",
		Data: &docPayload{Message: "hello"},
	}
	td.Init(context.Background(), &td.Data, func() { td.Data = nil }, coll, k)

	require.NoError(t, td.Create())

	require.NoError(t, td.Update(func() bool {
		td.Data.Message = "world"
		return true
	}))

	loaded := &testDoc{
		ID:   "10000",
		Data: &docPayload{},
	}
	loaded.Init(context.Background(), &loaded.Data, func() { loaded.Data = nil }, coll, k)
	require.NoError(t, loaded.Load())
	require.Equal(t, "world", loaded.Data.Message)

	require.NoError(t, loaded.Delete())
}

func TestDocumentBase_HashCacheReadThrough(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("doc-cache")
	require.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "cache-1")
	require.NoError(t, err)

	cache := newMemoryHashCache()
	td := &testDoc{ID: "cache-1", Data: &docPayload{Message: "cached"}}
	td.InitWithCache(context.Background(), &td.Data, func() { td.Data = nil }, coll, k, cache)
	require.NoError(t, td.Create())

	// Delete underlying DB document; Load should still succeed from HASH cache.
	require.NoError(t, coll.Delete(context.Background(), k))

	loaded := &testDoc{ID: "cache-1", Data: &docPayload{}}
	loaded.InitWithCache(context.Background(), &loaded.Data, func() { loaded.Data = nil }, coll, k, cache)
	require.NoError(t, loaded.Load())
	require.Equal(t, "cached", loaded.Data.Message)
	require.Equal(t, noptions.Version(1), loaded.version)
}

func TestDocumentBase_SaveAsyncWithoutWriteBackPersists(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("doc-sync-fallback")
	require.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "sync-1")
	require.NoError(t, err)

	td := &testDoc{ID: "sync-1", Data: &docPayload{Message: "hello"}}
	td.Init(context.Background(), &td.Data, func() { td.Data = nil }, coll, k)
	require.NoError(t, td.Create())

	require.NoError(t, td.UpdateAsync(func() bool {
		td.Data.Message = "persisted"
		return true
	}))

	var retrieved docPayload
	ver, err := coll.Get(context.Background(), k, noptions.WithDestination(&retrieved))
	require.NoError(t, err)
	require.Equal(t, "persisted", retrieved.Message)
	require.Equal(t, noptions.Version(2), ver)
}

type failPublishMQ struct {
	handler miface.SubResponseHandler
}

func (m *failPublishMQ) Publish(topic string, opts ...miface.PubOption) error {
	return errors.New("publish failed")
}

func (m *failPublishMQ) Subscribe(
	ctx context.Context,
	topic string,
	handler miface.SubResponseHandler,
	opts ...miface.SubOption,
) (miface.Subscription, error) {
	m.handler = handler
	return &noopSub{}, nil
}

func TestDocumentBase_SaveAsyncPublishFailureFallsBack(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("doc-publish-fail")
	require.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "pub-fail-1")
	require.NoError(t, err)

	cache := newMemoryHashCache()
	mq := &failPublishMQ{}
	td := &testDoc{ID: "pub-fail-1", Data: &docPayload{Message: "hello"}}
	td.InitWithCache(context.Background(), &td.Data, func() { td.Data = nil }, coll, k, cache)
	require.NoError(t, td.Create())
	require.NoError(t, td.EnableWriteBackWithMQ(mq, time.Millisecond))

	require.NoError(t, td.UpdateAsync(func() bool {
		td.Data.Message = "fallback"
		return true
	}))

	require.Eventually(t, func() bool {
		var got docPayload
		_, err := coll.Get(context.Background(), k, noptions.WithDestination(&got))
		return err == nil && got.Message == "fallback"
	}, time.Second, 10*time.Millisecond)
}

func TestDocumentBase_SaveAsyncDisableWriteBackNoPanic(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("doc-disable-wb")
	require.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "disable-1")
	require.NoError(t, err)

	cache := newMemoryHashCache()
	mq := &capturingMQ{}
	td := &testDoc{ID: "disable-1", Data: &docPayload{Message: "hello"}}
	td.InitWithCache(context.Background(), &td.Data, func() { td.Data = nil }, coll, k, cache)
	require.NoError(t, td.Create())
	require.NoError(t, td.EnableWriteBackWithMQ(mq, 30*time.Millisecond))

	require.NoError(t, td.UpdateAsync(func() bool {
		td.Data.Message = "async"
		return true
	}))
	td.DisableWriteBack()

	require.Eventually(t, func() bool {
		mq.mu.Lock()
		defer mq.mu.Unlock()
		return len(mq.messages) >= 1
	}, time.Second, 10*time.Millisecond)
}

func TestDocumentBase_UpdateAsyncWriteBack(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("doc-writeback")
	require.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "wb-1")
	require.NoError(t, err)

	cache := newMemoryHashCache()
	mq := &capturingMQ{}
	worker := NewWriteBackWorker(mq, provider, logger).WithCache(cache)
	require.NoError(t, worker.Start())
	defer worker.Stop()

	td := &testDoc{ID: "wb-1", Data: &docPayload{Message: "hello"}}
	td.InitWithCache(context.Background(), &td.Data, func() { td.Data = nil }, coll, k, cache)
	require.NoError(t, td.Create())
	require.Equal(t, int64(1), td.epoch)
	require.NoError(t, td.EnableWriteBackWithMQ(mq, time.Millisecond))

	require.NoError(t, td.UpdateAsync(func() bool {
		td.Data.Message = "async"
		return true
	}))
	require.Equal(t, noptions.Version(2), td.version)

	require.Eventually(t, func() bool {
		return worker.GetMetrics().ProcessedCount >= 1
	}, 2*time.Second, 20*time.Millisecond)

	var retrieved docPayload
	ver, err := coll.Get(context.Background(), k, noptions.WithDestination(&retrieved))
	require.NoError(t, err)
	require.Equal(t, "async", retrieved.Message)
	require.Equal(t, noptions.Version(2), ver)
}

func TestDocumentBase_DeleteRecreateFencesOldWriteBack(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("doc-epoch-fence")
	require.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "epoch-fence-1")
	require.NoError(t, err)

	cache := newMemoryHashCache()
	mq := &capturingMQ{}
	worker := NewWriteBackWorker(mq, provider, logger).WithCache(cache)
	require.NoError(t, worker.Start())
	defer worker.Stop()

	td := &testDoc{ID: "epoch-fence-1", Data: &docPayload{Message: "gen1"}}
	td.InitWithCache(context.Background(), &td.Data, func() { td.Data = nil }, coll, k, cache)
	require.NoError(t, td.Create())
	require.NoError(t, td.EnableWriteBackWithMQ(mq, 80*time.Millisecond))

	require.NoError(t, td.UpdateAsync(func() bool {
		td.Data.Message = "stale-gen1"
		return true
	}))

	// Recreate under a new epoch before the delayed write-back lands.
	require.NoError(t, td.Delete())
	td.Data = &docPayload{Message: "gen2"}
	require.NoError(t, td.Create())
	require.Equal(t, int64(2), td.epoch)

	require.Eventually(t, func() bool {
		return worker.GetMetrics().FailedCount >= 1
	}, 2*time.Second, 20*time.Millisecond)

	var got docPayload
	_, err = coll.Get(context.Background(), k, noptions.WithDestination(&got))
	require.NoError(t, err)
	require.Equal(t, "gen2", got.Message)
	require.Equal(t, int64(0), worker.GetMetrics().ProcessedCount)
}

func TestDocumentBase_ConcurrentUpdates(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	coll, err := provider.OpenDbDriver("doc-cas")
	require.NoError(t, err)

	k, err := key.NewKeyFromParts("demo", "cas-1")
	require.NoError(t, err)

	seed := &testDoc{ID: "cas-1", Data: &docPayload{Message: "seed"}}
	seed.Init(context.Background(), &seed.Data, func() { seed.Data = nil }, coll, k)
	require.NoError(t, seed.Create())

	const workers = 10
	const outerAttempts = 32

	var success atomic.Int32
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(idx int) {
			defer wg.Done()
			doc := &testDoc{ID: "cas-1", Data: &docPayload{}}
			doc.Init(context.Background(), &doc.Data, func() { doc.Data = nil }, coll, k)
			for attempt := 0; attempt < outerAttempts; attempt++ {
				if err := doc.Load(); err != nil {
					time.Sleep(time.Millisecond)
					continue
				}
				err := doc.Update(func() bool {
					doc.Data.Message = fmt.Sprintf("world+%d", idx)
					return true
				})
				if err == nil {
					success.Add(1)
					return
				}
			}
		}(i)
	}
	wg.Wait()

	require.Equal(t, int32(workers), success.Load(), "every worker should eventually CAS-update")

	var retrieved docPayload
	ver, err := coll.Get(context.Background(), k, noptions.WithDestination(&retrieved))
	require.NoError(t, err)
	require.Equal(t, noptions.Version(workers+1), ver, "expected version after create + %d updates", workers)
}
