package internal

import (
	"github.com/gstones/platform/services/common/nosql/document"
	"github.com/gstones/platform/services/common/nosql/errors"
	"sync"

	"github.com/snichols/gocb"
	"go.uber.org/zap"
)

// DocumentStoreProvider provides nosql.DocumentStore instances backed by Couchbase.
type DocumentStoreProvider struct {
	mutex   sync.Mutex
	cluster *gocb.Cluster
	buckets map[string]*DocumentStore
	logger  *zap.Logger
}

// NewDocumentStoreProvider constructs and returns a new DocumentStoreProvider backed by a Couchbase cluster.
func NewDocumentStoreProvider(connstr string, username string, password string, l *zap.Logger) (document.DocumentStoreProvider, error) {
	if c, err := gocb.Connect(connstr); err != nil {
		return nil, err
	} else {
		if username != "" && password != "" {
			c.Authenticate(gocb.PasswordAuthenticator{
				Username: username,
				Password: password,
			})
		}

		return &DocumentStoreProvider{
			cluster: c,
			buckets: map[string]*DocumentStore{},
			logger:  l,
		}, nil
	}
}

// OpenDocumentStore returns the DocumentStore with the given name.
func (d *DocumentStoreProvider) OpenDocumentStore(name string) (document.DocumentStore, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if s, ok := d.buckets[name]; ok {
		return s, nil
	} else if s, err := NewDocumentStore(name, d); err != nil {
		return nil, err
	} else {
		d.buckets[name] = s
		return s, nil
	}
}

// Shutdown closes the couchbase connection.
func (d *DocumentStoreProvider) Shutdown() error {
	//return d.cluster.Close()
	return nil
}

// AddStartingDocuments adds documents to the underlying document stores in bulk.  This is not supported
// in this driver.
func (d *DocumentStoreProvider) AddStartingDocuments(documents []document.StartingDocument) error {
	return errors.ErrNotImplemented
}
