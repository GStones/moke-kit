package nosql

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/fxmain"
	"github.com/gstones/moke-kit/mq/pkg/mfx"
	"github.com/gstones/moke-kit/orm/nosql/cache"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/pkg/ofx"
)

type SubData struct {
	SubMessage string
	SubList    []string
}

// MarshalBinary 实现 encoding.BinaryMarshaler 接口
func (td *SubData) MarshalBinary() ([]byte, error) {
	return json.Marshal(td)
}

// UnmarshalBinary 实现 encoding.BinaryUnmarshaler 接口
func (td *SubData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, td)
}

type TestData struct {
	DocumentBase `json:"-" bson:"-"`
	ID           string
	Message      string
	AList        []string
	BMap         map[string]string
	SubData      *SubData
}

// MarshalBinary 实现 encoding.BinaryMarshaler 接口
func (td *TestData) MarshalBinary() ([]byte, error) {
	return json.Marshal(td)
}

// UnmarshalBinary 实现 encoding.BinaryUnmarshaler 接口
func (td *TestData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, td)
}

func createTestData(
	logger *zap.Logger,
	id string,
	mongoCollect diface.ICollection,
	client *redis.Client,
) *TestData {
	td := &TestData{
		ID:      id,
		Message: "hello",
		AList:   []string{"a", "b"},
		BMap:    map[string]string{"key1": "value1", "key2": "value2"},
		SubData: &SubData{
			SubMessage: "sub hello",
			SubList:    []string{},
		},
	}
	k, err := key.NewKeyFromParts("demo", td.ID)
	if err != nil {
		panic(err)
	}
	if client == nil {
		td.Init(context.Background(), &td, td.clear, mongoCollect, k)
	} else {
		td.InitWithCache(context.Background(), &td, td.clear, mongoCollect, k, cache.CreateRedisCache(logger, client))
	}
	return td
}

func (td *TestData) clear() {

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
		td := createTestData(log, "10000", mongoCollect, nil)
		err = td.Create()
		assert.NoError(t, err, "Failed to create document")
		assert.Equal(t, int64(1), td.version, "Initial version should be 1")

		// Test Update
		err = td.Update(func() bool {
			td.Message = "world"
			return true
		})
		assert.NoError(t, err, "Failed to update document")
		assert.Equal(t, int64(2), td.version, "Version should be 2 after update")

		// Test Load
		td.Load()
		assert.NoError(t, err, "Failed to load document")
		assert.Equal(t, "world", td.Message, "Loaded message doesn't match")

		// Test Delete
		err = td.Delete()
		assert.NoError(t, err, "Failed to delete document")

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
		tdInit := createTestData(log, "10000", mongoCollect, nil)
		err = tdInit.Create()
		assert.NoError(t, err, "Failed to create initial document")

		// Create multiple test data instances
		const numGoroutines = 10
		testDatas := make([]*TestData, numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			td := createTestData(log, "10000", mongoCollect, nil)
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
					data.Message = fmt.Sprintf("world+%d", index)
					return true
				})
				assert.NoError(t, err, "Concurrent update failed")
			}(i, td)
		}
		wg.Wait()

		// Verify final state
		finalTd := createTestData(log, "10000", mongoCollect, nil)
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

func TestDocumentBase_SaveAsync(t *testing.T) {
	var testModule = fx.Invoke(func(
		log *zap.Logger,
		natsParams mfx.MessageQueueParams,
		dProvider ofx.DocumentStoreParams,
		params ofx.RedisParams,
	) error {
		if err := NewWriteBackWorker(natsParams.MessageQueue, dProvider.DriverProvider, log).Start(); err != nil {
			t.Fatal(err)
		}

		mongoCollect, err := dProvider.DriverProvider.OpenDbDriver("wirteback")
		if err != nil {
			t.Fatal(err)
		}

		td := createTestData(log, "10000", mongoCollect, params.Cache)
		err = td.Load()
		assert.NoError(t, err, "Failed to create document")

		if err := td.EnableWriteBackWithMQ(natsParams.MessageQueue, time.Duration(-1)); err != nil {
			t.Fatal(err)
		}
		err = td.UpdateAsync(func() bool {
			td.Message = "msg"
			td.AList = append(td.AList, "c")
			td.SubData.SubList = append(td.SubData.SubList, "sub e")
			td.BMap["key3"] = "value3"
			td.SubData.SubMessage = "sub msg2"
			return true
		})
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(20 * time.Second)

		return nil
	})

	fxmain.Main(mfx.NatsModule, testModule)
}
