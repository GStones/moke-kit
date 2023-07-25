package document

import (
	"github.com/gstones/platform/services/common/nosql/errors"
	"math/rand"
	"reflect"
	"time"

	"github.com/gstones/platform/services/common/jsonx"
)

const MaxRetries = 5

// Document provides a generic JSON object that's loaded from a DocumentStore.
type Document struct {
	clear    func()
	dataType reflect.Type
	data     interface{}
	version  Version

	DocumentStore DocumentStore
	Key           Key
	cache         DocumentCache
}

// Init performs an in-place initialization of a Document.
func (d *Document) Init(data interface{}, clear func(), store DocumentStore, key Key) {
	d.clear = clear
	d.dataType = reflect.TypeOf(data)
	d.data = data
	d.DocumentStore = store
	d.Key = key
	d.cache = NewDocumentCache()
}
func (d *Document) InitWithCache(
	data interface{},
	clear func(),
	store DocumentStore,
	key Key,
	cache DocumentCache,
) {
	d.Init(data, clear, store, key)
	d.cache = cache

}

func (d *Document) InitWithVersion(data interface{}, clear func(), store DocumentStore, key Key, version Version) {
	d.Init(data, clear, store, key)
	d.version = version
}

// LoadFromString loads this entity from the provided string.
func (d *Document) LoadFromString(s string) (err error) {
	return jsonx.ParseString(s, d.data)
}

// Clear clears all data on this Document.
func (d *Document) Clear() {
	d.version = NoVersion
	d.clear()
}

// Create creates this Document if it doesn't already exist.
func (d *Document) Create() (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		WithSource(d.data),
	)
	return
}

// CreateExpiry creates this Document with an expiration value, if it doesn't already exist.
func (d *Document) CreateExpiry(expiry time.Duration) (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		WithSource(d.data),
		WithTTL(expiry),
	)
	return
}

// Load loads this Document from its store if it exists.
func (d *Document) Load() (err error) {
	d.clear()
	if ok := d.cache.GetCache(d.Key, d.data); !ok {
		if d.version, err = d.DocumentStore.Get(
			d.Key,
			WithDestination(d.data),
		); err != nil {
			return err
		} else {
			d.cache.SetCache(d.Key, d.data)
		}
	}
	return
}

// LoadAndTouch loads this Document from its store if it exists while
// simultaneously updating any Expiry setting.
func (d *Document) LoadAndTouch(expiry time.Duration) (err error) {
	d.clear()

	d.version, err = d.DocumentStore.Get(
		d.Key,
		WithDestination(d.data),
		WithTTL(expiry),
	)
	return
}

// Save saves this Document to the database if it's based on the latest version that Couchbase knows about.
func (d *Document) Save() (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		WithSource(d.data),
		WithVersion(d.version),
	)
	return
}

// SaveExpiry saves this Document to the database, with a new expiration value,
// if it's based on the latest version that Couchbase knows about.
func (d *Document) SaveExpiry(expiry time.Duration) (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		WithSource(d.data),
		WithVersion(d.version),
		WithTTL(expiry),
	)
	return
}

// doUpdate executes the provided function until it returns true and
// the Document can be stored using the provided update function in the database without issue.
// If too many attempts are made then the function fails with ErrTooManyRetries.
func (d *Document) doUpdate(f func() bool, u func() error) error {
	for r := 0; r < MaxRetries; r++ {
		if f() {
			if err := u(); err == nil {
				return nil
			} else {
				// SNICHOLS: This random sleep is to help with updating the same Entity multiple times.  See
				// https://github.com/89trillion/platform/services/issues/136 for more details.
				time.Sleep(time.Millisecond * time.Duration(rand.Float32()*float32(r+1)*5))
				if err := d.Load(); err != nil {
					// a failed load is a real error
					return err
				}
			}
		} else {
			return errors.ErrUpdateLogicFailed
		}
	}

	return errors.ErrTooManyRetries
}

// Invokes doUpdate with Save() as the update function.
func (d *Document) Update(f func() bool) error {
	if err := d.doUpdate(f, func() error {
		return d.Save()
	}); err != nil {
		return err
	} else {
		d.cache.DeleteCache(d.Key)
		return nil
	}

}

// Invokes doUpdate with SaveExpiry() as the update function.
func (d *Document) UpdateExpiry(f func() bool, expiry time.Duration) error {
	return d.doUpdate(f, func() error {
		return d.SaveExpiry(expiry)
	})

}

// Invokes document sub-mutation to add entry to end of array.
func (d *Document) PushBack(key string, path string, item interface{}) error {
	return d.DocumentStore.PushBack(Key{value: key}, path, item)
}

// Invokes document sub-mutation to add entry to end of array.
func (d *Document) Incr(key string, path string, item interface{}) error {
	return d.DocumentStore.PushBack(Key{value: key}, path, item)
}

// Remove removes the Document from the database.
func (d *Document) Remove() error {
	return d.DocumentStore.Remove(d.Key, WithAnyVersion())
}
