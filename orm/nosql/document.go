package nosql

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

const MaxRetries = 5

type DocumentBase struct {
	Key  key.Key
	Name string

	clear    func()
	dataType reflect.Type
	data     any
	version  noptions.Version

	DocumentStore diface.ICollection
	cache         diface.IDocumentCache
}

// Init performs an in-place initialization of a DocumentBase.
func (d *DocumentBase) Init(
	data any,
	clear func(),
	store diface.ICollection,
	key key.Key,
) {
	d.clear = clear
	d.dataType = reflect.TypeOf(data)
	d.data = data
	d.DocumentStore = store
	d.Key = key
	d.cache = diface.DefaultDocumentCache()
}

func (d *DocumentBase) InitWithCache(
	data any,
	clear func(),
	store diface.ICollection,
	key key.Key,
	cache diface.IDocumentCache,
) {
	d.Init(data, clear, store, key)
	d.cache = cache
}

func (d *DocumentBase) InitWithVersion(
	data any,
	clear func(),
	store diface.ICollection,
	key key.Key,
	version noptions.Version,
) {
	d.Init(data, clear, store, key)
	d.version = version
}

// Clear clears all data on this DocumentBase.
func (d *DocumentBase) Clear() {
	d.version = noptions.NoVersion
	d.clear()
}

// Create creates this DocumentBase if it doesn't already exist.
func (d *DocumentBase) Create() (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		noptions.WithSource(d.data),
	)
	return
}

// CreateExpiry creates this DocumentBase with an expiration value, if it doesn't already exist.
func (d *DocumentBase) CreateExpiry(expiry time.Duration) (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		noptions.WithSource(d.data),
		noptions.WithTTL(expiry),
	)
	return
}

// Load loads this DocumentBase from its store if it exists.
func (d *DocumentBase) Load() (err error) {
	d.clear()
	if ok := d.cache.GetCache(d.Key, d.data); !ok {
		if d.version, err = d.DocumentStore.Get(
			d.Key,
			noptions.WithDestination(d.data),
		); err != nil {
			return err
		} else {
			d.cache.SetCache(d.Key, d.data)
		}
	}
	return
}

// LoadAndTouch loads this DocumentBase from its store if it exists while
// simultaneously updating any Expiry setting.
func (d *DocumentBase) LoadAndTouch(expiry time.Duration) (err error) {
	d.clear()

	d.version, err = d.DocumentStore.Get(
		d.Key,
		noptions.WithDestination(d.data),
		noptions.WithTTL(expiry),
	)
	return
}

// Save saves this DocumentBase to the database if it's based on the latest version that Couchbase knows about.
func (d *DocumentBase) Save() (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		noptions.WithSource(d.data),
		noptions.WithVersion(d.version),
	)
	return
}

// SaveExpiry saves this DocumentBase to the database, with a new expiration value,
// if it's based on the latest version that Couchbase knows about.
func (d *DocumentBase) SaveExpiry(expiry time.Duration) (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		noptions.WithSource(d.data),
		noptions.WithVersion(d.version),
		noptions.WithTTL(expiry),
	)
	return
}

func (d *DocumentBase) doUpdate(f func() bool, u func() error) error {
	for r := 0; r < MaxRetries; r++ {
		if f() {
			if err := u(); err == nil {
				return nil
			} else {
				time.Sleep(time.Millisecond * time.Duration(rand.Float32()*float32(r+1)*5))
				if err := d.Load(); err != nil {
					// a failed load is a real error
					return err
				}
			}
		} else {
			return nerrors.ErrUpdateLogicFailed
		}
	}

	return nerrors.ErrTooManyRetries
}

func (d *DocumentBase) Update(f func() bool) error {
	if err := d.doUpdate(f, func() error {
		return d.Save()
	}); err != nil {
		return err
	} else {
		d.cache.DeleteCache(d.Key)
		return nil
	}
}

func (d *DocumentBase) UpdateExpiry(f func() bool, expiry time.Duration) error {
	return d.doUpdate(f, func() error {
		return d.SaveExpiry(expiry)
	})
}
