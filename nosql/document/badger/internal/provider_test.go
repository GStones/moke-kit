package internal

import (
	"github.com/gstones/platform/services/common/nosql/document"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/utils"
)

func TestDocumentStoreProvider(t *testing.T) {
	if tmp, err := utils.NewTempDir("provider_test"); err != nil {
		t.Fatal(err)
	} else {
		defer tmp.Cleanup()

		if provider := NewDocumentStoreProvider(tmp.Path(), 5*time.Minute, zap.NewNop()); provider == nil {
			t.Fatal("creating document store provider failed")
		} else {
			defer provider.Shutdown()

			if store, err := provider.OpenDocumentStore("test"); err != nil {
				t.Fatal(err)
			} else {
				key := document.NewKey(`/local/test`)
				if _, err := store.Set(key, document.WithSource("test")); err != nil {
					t.Fatal(err)
				} else if ok, err := store.Contains(key); err != nil {
					t.Fatal(err)
				} else if !ok {
					t.Fatal("boom")
				}
			}
		}
	}
}
