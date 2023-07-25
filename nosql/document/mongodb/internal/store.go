package internal

import (
	"context"
	"github.com/gstones/platform/services/common/dupblock"
	"github.com/gstones/platform/services/common/nosql/document"
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
)

// DocumentStore adapts a MongoDB Collection to the nosql.DocumentStore interface.
type DocumentStore struct {
	mutex          sync.Mutex
	name           string
	owner          *DocumentStoreProvider
	database       *mongo.Database
	bucketWatchers sync.Map
	logger         *zap.Logger
	collects       sync.Map
}

func (d *DocumentStore) getCollection(key document.Key) *mongo.Collection {
	name := key.Prefix()
	collect, _ := d.collects.LoadOrStore(name, d.database.Collection(name))
	return collect.(*mongo.Collection)
}

func (d *DocumentStore) Remove(key document.Key, opts ...document.Option) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.owner.timeout)
	defer cancel()
	collect := d.getCollection(key)
	_, err := collect.DeleteOne(ctx, bson.M{"_id": key.String()})
	if err != nil {
		return convertMongodbError(err, key.String())
	}
	return nil
}

// ApplyDUPBlock applies a document update block to the DocumentStore.
func (d *DocumentStore) ApplyDUPBlock(reader dupblock.Reader) error {
	panic("implement me")
}

func (d *DocumentStore) SetField(key document.Key, path string, opts ...document.Option) (document.Version, error) {
	panic("implement me")
}

func (d *DocumentStore) AddWatcher(key document.Key, watcher document.DocumentWatcher) {
	// 创建 Change Stream
	collect := d.getCollection(key)
	watchers, _ := d.bucketWatchers.LoadOrStore(key, NewVBucketWatcher(key, collect, d.logger))
	watchers.(*VBucketWatcher).AddWatcher(key, watcher)
}

func (d *DocumentStore) RemoveWatcher(key document.Key, watcher document.DocumentWatcher) {
	watchers, ok := d.bucketWatchers.Load(key)
	if !ok {
		return
	}
	if watchers.(*VBucketWatcher).StopWatchingDocument(watcher) {
		d.bucketWatchers.Delete(key)
	}
}

func (d *DocumentStore) WatchersForKey(key document.Key) []document.DocumentWatcher {
	//TODO implement me
	panic("implement me")
}

func (d *DocumentStore) PushBack(key document.Key, path string, item interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.owner.timeout)
	defer cancel()
	collect := d.getCollection(key)
	filter := bson.M{"_id": key.String()}
	update := bson.M{"$push": bson.M{path: item}}
	opt := options.Update().SetUpsert(true)
	if _, err := collect.UpdateOne(ctx, filter, update, opt); err != nil {
		return convertMongodbError(err, key.String())
	}

	return nil
}

func (d *DocumentStore) Touch(key document.Key, opts ...document.Option) (document.Version, error) {
	//TODO implement me

	panic("implement me")
}

func NewDocumentStore(name string, owner *DocumentStoreProvider) (*DocumentStore, error) {
	collection := owner.cluster.Database(name)

	return &DocumentStore{
		name:     name,
		owner:    owner,
		logger:   owner.logger,
		database: collection,
	}, nil
}

func (d *DocumentStore) Name() string {
	return d.name
}

// Contains checks to see if a document with the given keys exists in the DocumentStore.
func (d *DocumentStore) Contains(key document.Key) (contains bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.owner.timeout)
	defer cancel()
	collect := d.getCollection(key)
	filter := bson.M{"_id": key.String()}
	count, err := collect.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (d *DocumentStore) ListKeys(prefix string, opts ...document.ScanOption) (list []document.Key, err error) {
	panic("implement me")
}

func (d *DocumentStore) Set(key document.Key, opts ...document.Option) (document.Version, error) {
	if o, e := document.NewOptions(opts...); e != nil {
		return document.NoVersion, e
	} else if o.Source == nil {
		return document.NoVersion, errors2.ErrSourceIsNil
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), d.owner.timeout)
		defer cancel()
		update := bson.M{"$set": o.Source}
		collect := d.getCollection(key)
		opt := options.Update().SetUpsert(true)
		if _, err := collect.UpdateByID(ctx, key.String(), update, opt); err != nil {
			return document.NoVersion, convertMongodbError(err, key.String())
		}
	}
	return document.NoVersion, nil
}

// Get loads an existing document from the document store and returns its cas.  If no such document exists then
// this function fails.
func (d *DocumentStore) Get(key document.Key, opts ...document.Option) (document.Version, error) {
	if o, e := document.NewOptions(opts...); e != nil {
		return document.NoVersion, e
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), d.owner.timeout)
		defer cancel()
		filter := bson.M{"_id": key.String()}
		collect := d.getCollection(key)
		if err := collect.FindOne(ctx, filter).Decode(o.Destination); err != nil {
			return document.NoVersion, convertMongodbError(err, key.String())
		}
	}
	return document.NoVersion, nil
}

func (d *DocumentStore) Scan(prefix string, destOpt document.Option, scanOpts ...document.ScanOption) (amt int, err error) {
	panic("implement me")
}

func (d *DocumentStore) Incr(key document.Key, path string, num int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.owner.timeout)
	defer cancel()
	collect := d.getCollection(key)
	filter := bson.M{"_id": key.String()}
	update := bson.M{"$inc": bson.M{path: num}}
	opt := options.Update().SetUpsert(true)
	if _, err := collect.UpdateOne(ctx, filter, update, opt); err != nil {
		return convertMongodbError(err, key.String())
	}
	return nil

}
