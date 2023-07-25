package internal

import (
	"context"
	"github.com/gstones/platform/services/common/nosql/document"
	"github.com/gstones/platform/services/common/nosql/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"time"

	"go.uber.org/zap"
)

// DocumentStoreProvider provides nosql.DocumentStore instances backed by Couchbase.
type DocumentStoreProvider struct {
	name    string
	timeout time.Duration
	mutex   sync.Mutex
	cluster *mongo.Client
	buckets map[string]*DocumentStore
	logger  *zap.Logger
}

// NewDocumentStoreProvider constructs and returns a new DocumentStoreProvider backed by a Couchbase cluster.
func NewDocumentStoreProvider(connstr string, username string, password string, l *zap.Logger) (document.DocumentStoreProvider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	credential := options.Credential{
		Username: username,
		Password: password,
	}
	cOptions := options.Client().ApplyURI(connstr)
	if username != "" && password != "" {
		cOptions.SetAuth(credential)
	}
	if client, err := mongo.Connect(
		ctx,
		cOptions,
	); err != nil {
		return nil, err
	} else {
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			return nil, err
		}
		return &DocumentStoreProvider{
			cluster: client,
			buckets: map[string]*DocumentStore{},
			logger:  l,
			name:    "game",
			timeout: time.Second * 10,
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
