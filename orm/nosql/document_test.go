package nosql

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/mock"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

type docPayload struct {
	Message string `bson:"message"`
}

type testDoc struct {
	DocumentBase
	ID   string
	Data *docPayload
}

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

	// Workers may contend past a single DocumentBase.Update's MaxRetries (5).
	// Each worker reloads and retries until success so the test asserts CAS
	// progress without flaking under -race stampede.
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
