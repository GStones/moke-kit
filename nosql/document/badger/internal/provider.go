package internal

import (
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/dgraph-io/badger"

	"go.uber.org/zap"
)

type DocumentStoreProvider struct {
	sync.Mutex
	dir        string
	stores     map[string]*DocumentStore
	logger     *zap.Logger
	gcInterval time.Duration
}

// NewDocumentStoreProvider constructs and returns a new DocumentStoreProvider backed by badgerdb.
func NewDocumentStoreProvider(dir string, gcInterval time.Duration, l *zap.Logger) *DocumentStoreProvider {
	return &DocumentStoreProvider{
		dir:        dir,
		stores:     map[string]*DocumentStore{},
		logger:     l,
		gcInterval: gcInterval,
	}

}

func (d *DocumentStoreProvider) OpenDocumentStore(name string) (document.DocumentStore, error) {
	d.Lock()
	defer d.Unlock()

	if store, ok := d.stores[name]; ok {
		return store, nil
	} else if store, err := NewDocumentStore(d.dir, name, d.gcInterval, d.logger); err != nil {
		return nil, err
	} else {
		d.stores[name] = store
		return store, nil
	}
}

func (d *DocumentStoreProvider) Shutdown() (err error) {
	d.Lock()
	defer d.Unlock()

	for _, s := range d.stores {
		if e := s.Close(); e != nil {
			err = e
		}
	}

	return
}

// AddStartingDocuments adds documents to the underlying document stores in bulk.  This is not supported
// in this driver.
func (d *DocumentStoreProvider) AddStartingDocuments(documents []document.StartingDocument) error {
	return errors.ErrNotImplemented
}

// Replication of NewDocumentStore() functionality, modified to return only badger DB.
func NewBadgerStore(dir string, name string) (*badger.DB, error) {
	d := path.Join(dir, name)
	opts := badger.DefaultOptions(d)
	// Per https://github.com/dgraph-io/badger/issues/476, must set Truncate true on Windows.
	if runtime.GOOS == "windows" {
		opts.Truncate = true
	}

	if err := os.MkdirAll(opts.Dir, os.ModePerm); err != nil {
		return nil, err
	} else if db, _, err := openDB(opts); err != nil {
		return nil, convertBadgerError(err, "")
	} else {
		return db, nil
	}
}
