package badger

import (
	"github.com/gstones/platform/services/common/nosql/document"
	"time"

	"github.com/dgraph-io/badger"

	"github.com/gstones/platform/services/common/nosql/document/badger/internal"
	"go.uber.org/zap"
)

func NewDocumentStoreProvider(dir string, gcInterval time.Duration, l *zap.Logger) (document.DocumentStoreProvider, error) {
	return internal.NewDocumentStoreProvider(dir, gcInterval, l), nil
}

func NewBadgerStore(dir string, name string) (*badger.DB, error) {
	return internal.NewBadgerStore(dir, name)
}
