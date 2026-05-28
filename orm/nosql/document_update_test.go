package nosql

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

type updatePayload struct {
	Message string `json:"message"`
}

type fakeAtomicCache struct {
	mu          sync.Mutex
	version     noptions.Version
	data        *updatePayload
	failCASOnce bool
}

func (c *fakeAtomicCache) GetCache(_ context.Context, _ key.Key, doc any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	vc, ok := doc.(*VersionCache)
	if !ok {
		return false
	}
	if v, ok := vc.Version.(*noptions.Version); ok {
		*v = c.version
	}
	if dst, ok := vc.Data.(*updatePayload); ok && c.data != nil {
		*dst = *c.data
	}
	return true
}

func (c *fakeAtomicCache) SetCache(_ context.Context, _ key.Key, _ any, _ time.Duration) {}

func (c *fakeAtomicCache) DeleteCache(_ context.Context, _ key.Key) {}

func (c *fakeAtomicCache) CompareAndSwapCache(
	_ context.Context,
	_ key.Key,
	expectedVersion noptions.Version,
	doc any,
	_ time.Duration,
) (noptions.Version, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.failCASOnce {
		c.failCASOnce = false
		return noptions.NoVersion, nerrors.ErrVersionNotMatch
	}
	if expectedVersion != c.version {
		return noptions.NoVersion, nerrors.ErrVersionNotMatch
	}
	vc := doc.(*VersionCache)
	nextVersion := vc.Version.(noptions.Version)
	if payload, ok := vc.Data.(*updatePayload); ok {
		copyPayload := *payload
		c.data = &copyPayload
	}
	c.version = nextVersion
	return nextVersion, nil
}

type fakeCollection struct {
	mu            sync.Mutex
	version       noptions.Version
	failStoreOnce bool
	data          *updatePayload
	successCh     chan struct{}
}

func (c *fakeCollection) GetName() string { return "fake" }

func (c *fakeCollection) Set(_ context.Context, _ key.Key, opts ...noptions.Option) (noptions.Version, error) {
	o, err := noptions.NewOptions(opts...)
	if err != nil {
		return noptions.NoVersion, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.failStoreOnce {
		c.failStoreOnce = false
		return noptions.NoVersion, nerrors.ErrVersionNotMatch
	}
	if o.Version != c.version {
		return noptions.NoVersion, nerrors.ErrVersionNotMatch
	}
	if payload, ok := o.Source.(*updatePayload); ok {
		copyPayload := *payload
		c.data = &copyPayload
	}
	c.version++
	select {
	case c.successCh <- struct{}{}:
	default:
	}
	return c.version, nil
}

func (c *fakeCollection) Get(_ context.Context, _ key.Key, _ ...noptions.Option) (noptions.Version, error) {
	return c.version, nil
}

func (c *fakeCollection) Delete(_ context.Context, _ key.Key) error { return nil }

func (c *fakeCollection) Incr(_ context.Context, _ key.Key, _ string, _ int32) (int64, error) {
	return 0, nil
}

func TestDocumentUpdate_AtomicCacheRetryAndAsyncWriteBack(t *testing.T) {
	k, err := key.NewKeyFromParts("test", "doc", "1")
	require.NoError(t, err)

	payload := &updatePayload{Message: "old"}
	cache := &fakeAtomicCache{version: 1, data: &updatePayload{Message: "old"}, failCASOnce: true}
	store := &fakeCollection{version: 1, failStoreOnce: true, successCh: make(chan struct{}, 1)}

	var doc DocumentBase
	doc.InitWithCache(context.Background(), payload, func() {}, store, k, cache)
	doc.version = 1

	attempts := 0
	err = doc.Update(func() bool {
		attempts++
		payload.Message = "new"
		return true
	})
	require.NoError(t, err)
	require.Equal(t, 2, attempts)
	require.Equal(t, noptions.Version(2), doc.version)

	select {
	case <-store.successCh:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for async write-back")
	}

	require.Equal(t, noptions.Version(2), store.version)
	require.NotNil(t, store.data)
	require.Equal(t, "new", store.data.Message)
}
