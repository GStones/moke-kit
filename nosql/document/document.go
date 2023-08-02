package document

import (
	"math/rand"
	"reflect"
	"time"

	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/document/key"
	"moke-kit/nosql/errors"
)

const MaxRetries = 5

type DocumentBase struct {
	Key      key.Key
	clear    func()
	dataType reflect.Type
	data     interface{}
	version  diface.Version

	DocumentStore diface.ICollection
	cache         diface.IDocumentCache
}

// Init performs an in-place initialization of a DocumentBase.
func (d *DocumentBase) Init(
	data interface{},
	clear func(),
	key key.Key,
	store diface.ICollection,

) {
	d.clear = clear
	d.dataType = reflect.TypeOf(data)
	d.data = data
	d.DocumentStore = store
	d.Key = key
	d.cache = diface.DefaultDocumentCache()
}

func (d *DocumentBase) InitWithCache(
	data interface{},
	clear func(),
	store diface.ICollection,
	key key.Key,
	cache diface.IDocumentCache,
) {
	d.Init(data, clear, key, store)
	d.cache = cache
}

func (d *DocumentBase) InitWithVersion(
	data interface{},
	clear func(),
	store diface.ICollection,
	key key.Key,
	version diface.Version,
) {
	d.Init(data, clear, key, store)
	d.version = version
}

// Clear clears all data on this DocumentBase.
func (d *DocumentBase) Clear() {
	d.version = diface.NoVersion
	d.clear()
}

// Create creates this DocumentBase if it doesn't already exist.
func (d *DocumentBase) Create() (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		diface.WithSource(d.data),
	)
	return
}

// CreateExpiry creates this DocumentBase with an expiration value, if it doesn't already exist.
func (d *DocumentBase) CreateExpiry(expiry time.Duration) (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		diface.WithSource(d.data),
		diface.WithTTL(expiry),
	)
	return
}

// Load loads this DocumentBase from its store if it exists.
func (d *DocumentBase) Load() (err error) {
	d.clear()
	if ok := d.cache.GetCache(d.Key, d.data); !ok {
		if d.version, err = d.DocumentStore.Get(
			d.Key,
			diface.WithDestination(d.data),
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
		diface.WithDestination(d.data),
		diface.WithTTL(expiry),
	)
	return
}

// Save saves this DocumentBase to the database if it's based on the latest version that Couchbase knows about.
func (d *DocumentBase) Save() (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		diface.WithSource(d.data),
		diface.WithVersion(d.version),
	)
	return
}

// SaveExpiry saves this DocumentBase to the database, with a new expiration value,
// if it's based on the latest version that Couchbase knows about.
func (d *DocumentBase) SaveExpiry(expiry time.Duration) (err error) {
	d.version, err = d.DocumentStore.Set(
		d.Key,
		diface.WithSource(d.data),
		diface.WithVersion(d.version),
		diface.WithTTL(expiry),
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
			return errors.ErrUpdateLogicFailed
		}
	}

	return errors.ErrTooManyRetries
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
