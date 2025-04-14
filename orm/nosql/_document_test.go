package nosql

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/fxmain"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/pkg/ofx"
)

type Data struct {
	Message string
}

type TestData struct {
	DocumentBase
	ID   string `bson:"_id"`
	Data *Data  `bson:"data"`
}

func createTestData(id string, mongoCollect diface.ICollection) *TestData {
	td := &TestData{
		ID: id,
		Data: &Data{
			Message: "hello",
		},
	}
	k, err := key.NewKeyFromParts("demo", td.ID)
	if err != nil {
		panic(err)
	}
	td.Init(context.Background(), &td.Data, td.clear, mongoCollect, k)
	return td
}

func (td *TestData) clear() {
	td.Data = nil
}

func TestDocument_CRUD(t *testing.T) {
	os.Setenv("DATABASE_URL", "mock://127.0.0.1:27017")

	var testModule = fx.Invoke(func(
		log *zap.Logger,
		dProvider ofx.DocumentStoreParams,
	) error {
		mongoCollect, err := dProvider.DriverProvider.OpenDbDriver("test")
		if err != nil {
			t.Fatal(err)
		}

		// Test Create
		td := createTestData("10000", mongoCollect)
		err = td.Create()
		assert.NoError(t, err, "Failed to create document")
		assert.Equal(t, int64(1), td.version, "Initial version should be 1")

		// Test Update
		err = td.Update(func() bool {
			td.Data.Message = "world"
			return true
		})
		assert.NoError(t, err, "Failed to update document")
		assert.Equal(t, int64(2), td.version, "Version should be 2 after update")

		// Test Load
		td.Load()
		assert.NoError(t, err, "Failed to load document")
		assert.Equal(t, "world", td.Data.Message, "Loaded message doesn't match")

		// Test Delete
		err = td.Delete()
		assert.NoError(t, err, "Failed to delete document")

		os.Exit(0)
		return nil
	})

	fxmain.Main(testModule)
}

func TestDocument_ConcurrentUpdates(t *testing.T) {
	os.Setenv("DATABASE_URL", "mock://127.0.0.1:27017")

	var concurrentModule = fx.Invoke(func(
		log *zap.Logger,
		dProvider ofx.DocumentStoreParams,
	) error {
		mongoCollect, err := dProvider.DriverProvider.OpenDbDriver("test_concurrent")
		if err != nil {
			t.Fatal(err)
		}

		// Initialize document
		tdInit := createTestData("10000", mongoCollect)
		err = tdInit.Create()
		assert.NoError(t, err, "Failed to create initial document")

		// Create multiple test data instances
		const numGoroutines = 10
		testDatas := make([]*TestData, numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			td := createTestData("10000", mongoCollect)
			err = td.Load()
			assert.NoError(t, err, "Failed to load document for concurrent test")
			testDatas[i] = td
		}

		// Perform concurrent updates
		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		for i, td := range testDatas {
			go func(index int, data *TestData) {
				defer wg.Done()
				err := data.Update(func() bool {
					data.Data.Message = fmt.Sprintf("world+%d", index)
					return true
				})
				assert.NoError(t, err, "Concurrent update failed")
			}(i, td)
		}
		wg.Wait()

		// Verify final state
		finalTd := createTestData("10000", mongoCollect)
		err = finalTd.Load()
		assert.NoError(t, err, "Failed to load final state")
		assert.Equal(t, int64(11), finalTd.version, "Final version should be 11")

		// Cleanup
		err = finalTd.Delete()
		assert.NoError(t, err, "Failed to cleanup test document")

		os.Exit(0)
		return nil
	})

	fxmain.Main(concurrentModule)
}
