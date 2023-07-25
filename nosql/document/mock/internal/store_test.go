package internal

import (
	"testing"

	"github.com/gstones/platform/services/common/nosql/pkg/tests/common_tests"
	"github.com/gstones/platform/services/common/utils"
)

func TestNewDocumentStoreProvider(t *testing.T) {
	if tmp, err := utils.NewTempDir("provider_test"); err != nil {
		t.Fatal(err)
	} else {
		defer tmp.Cleanup()
		testStoreName := "testStore"
		if provider, err := NewDocumentStoreProvider(); err != nil {
			t.Fatal("Error encountered creating a new document store provider:", err)
		} else {
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
}
