package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"reflect"
	"time"

	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/document/key"
	"moke-kit/nosql/nerrors"
)

const MaxRetries = 5

type MongoBase struct {
	Key      key.Key
	clear    func()
	dataType reflect.Type
	data     interface{}
	version  diface.Version

	collection *mongo.Collection
	cache      diface.IDocumentCache
}

// Init performs an in-place initialization of a MongoBase.
func (d *MongoBase) Init(
	data interface{},
	clear func(),
	key key.Key,
	db *mongo.Database,
) {
	cOpts := &options.CollectionOptions{}
	cOpts.SetBSONOptions(&options.BSONOptions{
		UseJSONStructTags: true,
	})
	d.clear = clear
	d.dataType = reflect.TypeOf(data)
	d.data = data
	d.collection = db.Collection(key.Prefix(), cOpts)
	d.Key = key
	d.cache = diface.DefaultDocumentCache()
}

func (d *MongoBase) InitWithCache(
	data interface{},
	clear func(),
	collection *mongo.Database,
	key key.Key,
	cache diface.IDocumentCache,
) {
	d.Init(data, clear, key, collection)
	d.cache = cache
}

// Clear clears all data on this MongoBase.
func (d *MongoBase) Clear() {
	d.clear()
}

// Create creates this MongoBase if it doesn't already exist.
func (d *MongoBase) Create() error {
	if _, err := d.collection.InsertOne(context.Background(), d.data); err != nil {
		return err
	}
	return nil
}

// Load loads this MongoBase from its store if it exists.
func (d *MongoBase) Load() (err error) {
	d.clear()
	if ok := d.cache.GetCache(d.Key, d.data); !ok {
		filter := bson.M{"_id": d.Key.String()}
		if res := d.collection.FindOne(context.Background(), filter); res.Err() != nil {
			if errors.Is(res.Err(), mongo.ErrNoDocuments) {
				return nerrors.ErrNotFound
			}
			return res.Err()
		} else {
			bData := &bson.Raw{}
			if err = res.Decode(bData); err != nil {
				return err
			}
			version := bData.Lookup("version")
			vInt := version.Int64()
			d.version = vInt
			data := bData.Lookup("data")
			if err = data.Unmarshal(d.data); err != nil {
				return err
			}
			d.cache.SetCache(d.Key, d.data)
		}
	}
	return
}

// Save saves this MongoBase to the database if it's based on the latest version that Couchbase knows about.
func (d *MongoBase) Save() error {
	if _, err := d.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": d.Key.String(), "version": d.version},
		d.data,
	); err != nil {
		return err
	}
	return nil
}

func (d *MongoBase) doUpdate(f func() bool, u func() error) error {
	for r := 0; r < MaxRetries; r++ {
		if f() {
			if err := u(); err == nil {
				return nil
			} else {
				time.Sleep(time.Millisecond * time.Duration(rand.Float32()*float32(r+1)*5))
				if err := d.Load(); err != nil {
					return err
				}
			}
		} else {
			return nerrors.ErrUpdateLogicFailed
		}
	}

	return nerrors.ErrTooManyRetries
}

func (d *MongoBase) Update(f func() bool) error {
	if err := d.doUpdate(f, func() error {
		return d.Save()
	}); err != nil {
		return err
	} else {
		d.cache.DeleteCache(d.Key)
		return nil
	}
}
