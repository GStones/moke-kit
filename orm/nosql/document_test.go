package nosql

import (
	"context"
	"fmt"
	"sync"
	"testing"

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

	const n = 10
	docs := make([]*testDoc, n)
	for i := 0; i < n; i++ {
		d := &testDoc{ID: "cas-1", Data: &docPayload{}}
		d.Init(context.Background(), &d.Data, func() { d.Data = nil }, coll, k)
		require.NoError(t, d.Load())
		docs[i] = d
	}

	var wg sync.WaitGroup
	wg.Add(n)
	errs := make(chan error, n)
	for i, d := range docs {
		go func(idx int, doc *testDoc) {
			defer wg.Done()
			errs <- doc.Update(func() bool {
				doc.Data.Message = fmt.Sprintf("world+%d", idx)
				return true
			})
		}(i, d)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	var retrieved docPayload
	ver, err := coll.Get(context.Background(), k, noptions.WithDestination(&retrieved))
	require.NoError(t, err)
	require.Equal(t, noptions.Version(n+1), ver, "expected version after create + %d updates", n)
}
