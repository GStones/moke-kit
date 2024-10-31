package nosql

import (
	"fmt"
	log2 "log"
	"sync"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/fxmain"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/pkg/ofx"
)

// TestDocumentBase_Update is a test for document base update.
// Tips: require mongo server running.
func TestDocumentBase_Update(t *testing.T) {
	fxmain.Main(
		TestModule,
	)
}

// TestDocumentBase_Update is a test for document base update.
// Tips: require mongo server running.
func TestConcurrentUpdate(t *testing.T) {
	fxmain.Main(
		ConcurrentUpdateModule,
	)
}

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
	td.Init(&td.Data, td.clear, mongoCollect, k)
	return td
}

func (td *TestData) clear() {
	td.Data = nil
}

var TestModule = fx.Invoke(func(
	log *zap.Logger,
	dProvider ofx.DocumentStoreParams,
) error {
	// init mongo collection
	mongoCollect, err := dProvider.DriverProvider.OpenDbDriver("test")
	if err != nil {
		panic(err)
	}
	td := createTestData("10000", mongoCollect)
	if err := td.Create(); err != nil { // create document version =1
		panic(err)
	} else if err := td.Update(func() bool { // update document version =2
		// update logic
		td.Data.Message = "world"
		return true
	}); err != nil {
		panic(err)
	}

	td2 := createTestData("10000", mongoCollect)
	if err := td2.Load(); err != nil {
		panic(err)
	}
	if td2.Data.Message != "world" {
		log.Panic("message is not expect:world", zap.String("msg", td2.Data.Message))
	} else if td2.version != 2 {
		log.Panic("version is not expect:2", zap.Int64("version", td2.version))
	} else if err := td2.Delete(); err != nil {
		panic(err)
	}
	log2.Fatal("test done")
	return nil
})

var ConcurrentUpdateModule = fx.Invoke(func(
	log *zap.Logger,
	dProvider ofx.DocumentStoreParams,
) error {
	// init mongo collection
	mongoCollect, err := dProvider.DriverProvider.OpenDbDriver("test_concurrent")
	if err != nil {
		panic(err)
	}
	tdInit := createTestData("10000", mongoCollect)
	if err := tdInit.Create(); err != nil {
		panic(err)
	}

	// create 10 document
	testDatas := make([]*TestData, 0)
	for i := 0; i < 10; i++ {
		td := createTestData("10000", mongoCollect)
		if err := td.Load(); err != nil { // create document version =1
			panic(err)
		} else {
			testDatas = append(testDatas, td)
		}
	}

	// concurrent update
	wg := sync.WaitGroup{}
	wg.Add(10)
	for k, td := range testDatas {
		go func(td *TestData) {
			log.Info("start update", zap.String("msg:", td.Data.Message))
			if err := td.Update(func() bool { // update document version =2
				// update logic
				td.Data.Message = fmt.Sprint("world+", k)
				return true
			}); err != nil {
				panic(err)
			}
			wg.Done()
		}(td)
	}
	wg.Wait()

	td := createTestData("10000", mongoCollect)
	if err := td.Load(); err != nil {
		panic(err)
	}
	if td.version != 11 {
		log.Panic("version is not expect:11", zap.Int64("version", td.version))
	} else if err := td.Delete(); err != nil {
		panic(err)
	}
	log2.Fatal("test done")

	return nil
})
