package orm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/mock"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
	"github.com/gstones/moke-kit/test/utils"
)

// TestData represents test document structure
type TestData struct {
	Message string `bson:"message,omitempty"`
	Count   int    `bson:"count,omitempty"`
}

// TestDocument represents a complete document for testing
type TestDocument struct {
	ID   string    `bson:"_id"`
	Data *TestData `bson:"data"`
}

// TestKey tests the key functionality
func TestKey(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("CreateKey", func(t *testing.T) {
		// Test creating key from parts
		k, err := key.NewKeyFromParts("namespace", "entity", "id123")
		helper.RequireNoError(err)
		helper.AssertNotNil(k)

		keyStr := k.String()
		helper.AssertTrue(len(keyStr) > 0, "Key string should not be empty")
	})

	t.Run("KeyParts", func(t *testing.T) {
		parts := []string{"test", "user", "12345"}
		k, err := key.NewKeyFromParts(parts...)
		helper.RequireNoError(err)

		// Test key properties
		helper.AssertNotNil(k)
		helper.AssertTrue(len(k.String()) > 0)
	})

	t.Run("EmptyParts", func(t *testing.T) {
		// Test with empty parts
		_, err := key.NewKeyFromParts()
		helper.AssertNotNil(err, "Should return error for empty parts")
	})
}

// TestMockDocumentProvider tests the mock document provider
func TestMockDocumentProvider(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)

	t.Run("OpenDbDriver", func(t *testing.T) {
		collectionName := "test-collection"

		collection, err := provider.OpenDbDriver(collectionName)
		helper.RequireNoError(err)
		helper.AssertNotNil(collection)

		// Test collection name
		name := collection.GetName()
		helper.AssertEqual(collectionName, name)
	})

	t.Run("Shutdown", func(t *testing.T) {
		err := provider.Shutdown()
		helper.AssertNoError(err)
	})
}

// TestMockCollection tests the mock collection functionality using the actual mock implementation
func TestMockCollection(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	logger := zaptest.NewLogger(t)
	provider := mock.NewMockDriverProvider(logger)
	collection, err := provider.OpenDbDriver("test-collection")
	helper.RequireNoError(err)

	t.Run("SetAndGet", func(t *testing.T) {
		k, err := key.NewKeyFromParts("test", "document", "123")
		helper.RequireNoError(err)

		data := &TestData{
			Message: "test message",
			Count:   42,
		}

		// Test Set
		version, err := collection.Set(helper.Context(), k, noptions.WithSource(data))
		helper.RequireNoError(err)
		helper.AssertTrue(version > 0)

		// Test Get
		var retrieved TestData
		version, err = collection.Get(helper.Context(), k, noptions.WithDestination(&retrieved))
		helper.RequireNoError(err)
		helper.AssertTrue(version > 0)
	})

	t.Run("Delete", func(t *testing.T) {
		k, err := key.NewKeyFromParts("test", "document", "delete-me")
		helper.RequireNoError(err)

		err = collection.Delete(helper.Context(), k)
		helper.AssertNoError(err)
	})

	t.Run("Increment", func(t *testing.T) {
		k, err := key.NewKeyFromParts("test", "counter", "increment-test")
		helper.RequireNoError(err)

		newValue, err := collection.Incr(helper.Context(), k, "count", 5)
		helper.RequireNoError(err)
		helper.AssertTrue(newValue >= 0)
	})
}

// TestNoSQLOptions tests the NoSQL options
func TestNoSQLOptions(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	t.Run("BasicOptions", func(t *testing.T) {
		// Test default options
		opts, err := noptions.NewOptions()
		helper.RequireNoError(err)
		helper.AssertNotNil(opts)
		helper.AssertEqual(noptions.NoVersion, opts.Version)

		// Test with version
		version := noptions.Version(123)
		opts, err = noptions.NewOptions(noptions.WithVersion(version))
		helper.RequireNoError(err)
		helper.AssertEqual(version, opts.Version)

		// Test with any version
		opts, err = noptions.NewOptions(noptions.WithAnyVersion())
		helper.RequireNoError(err)
		helper.AssertTrue(opts.AnyVersion)

		// Test with TTL
		ttl := 5 * time.Minute
		opts, err = noptions.NewOptions(noptions.WithTTL(ttl))
		helper.RequireNoError(err)
		helper.AssertEqual(ttl, opts.TTL)
	})

	t.Run("SourceDestinationOptions", func(t *testing.T) {
		data := &TestData{Message: "test", Count: 42}

		// Test with source
		opts, err := noptions.NewOptions(noptions.WithSource(data))
		helper.RequireNoError(err)
		helper.AssertEqual(data, opts.Source)

		// Test with destination
		var dest TestData
		opts, err = noptions.NewOptions(noptions.WithDestination(&dest))
		helper.RequireNoError(err)
		helper.AssertEqual(&dest, opts.Destination)
	})

	t.Run("ConflictingOptions", func(t *testing.T) {
		// Test conflicting version options
		_, err := noptions.NewOptions(
			noptions.WithVersion(123),
			noptions.WithAnyVersion(),
		)
		helper.AssertNotNil(err, "Should return error for conflicting version options")
	})
}

// TestDocumentCache tests the document cache functionality
func TestDocumentCache(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	cache := diface.DefaultDocumentCache()

	t.Run("CacheOperations", func(t *testing.T) {
		k, err := key.NewKeyFromParts("test", "cache", "item1")
		helper.RequireNoError(err)

		data := &TestData{Message: "cached data", Count: 100}

		// Test cache set (default implementation does nothing)
		cache.SetCache(helper.Context(), k, data, 5*time.Minute)

		// Test cache get (default implementation returns false)
		var retrieved TestData
		found := cache.GetCache(helper.Context(), k, &retrieved)
		helper.AssertFalse(found, "Default cache should not find items")

		// Test cache delete (default implementation does nothing)
		cache.DeleteCache(helper.Context(), k)
	})
}

// BenchmarkMockCollection benchmarks the mock collection operations
func BenchmarkMockCollection(b *testing.B) {
	logger := zaptest.NewLogger(&testing.T{})
	provider := mock.NewMockDriverProvider(logger)
	collection, err := provider.OpenDbDriver("benchmark-collection")
	require.NoError(b, err)

	data := &TestData{
		Message: "benchmark message",
		Count:   42,
	}

	ctx := context.Background()

	b.Run("Set", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				k, err := key.NewKeyFromParts("benchmark", "set", "test")
				if err != nil {
					b.Error(err)
					continue
				}

				_, err = collection.Set(ctx, k, noptions.WithSource(data))
				if err != nil {
					b.Error(err)
				}
			}
		})
	})

	b.Run("Get", func(b *testing.B) {
		k, err := key.NewKeyFromParts("benchmark", "get", "test")
		require.NoError(b, err)

		var retrieved TestData

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := collection.Get(ctx, k, noptions.WithDestination(&retrieved))
				if err != nil {
					b.Error(err)
				}
			}
		})
	})
}
