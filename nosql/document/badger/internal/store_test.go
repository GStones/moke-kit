package internal

import (
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/nosql/pkg/tests/common_tests"
	"github.com/gstones/platform/services/common/utils"
)

func TestNewDocumentStoreProvider(t *testing.T) {
	if tmp, err := utils.NewTempDir("provider_test"); err != nil {
		t.Fatal(err)
	} else {
		defer tmp.Cleanup()
		testStoreName := "testStore"
		provider := NewDocumentStoreProvider(tmp.Path(), 5*time.Minute, zap.NewNop())
		if store, err := provider.OpenDocumentStore(testStoreName); err != nil {
			t.Fatal("Error opening document store", testStoreName, ":", err)
		} else if store.Name() != testStoreName {
			t.Fatal("Created store does not use the provided document store name.")
		} else if err := common_tests.StoreCommonTest(store); err != nil {
			t.Fatal(err)
		} else if err := provider.Shutdown(); err != nil {
			t.Fatal("Error encountered shutting down the document store provider:", err)
		}
	}
}
