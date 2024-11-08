package nosql

import (
	"math/rand"
	"time"

	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

const MaxRetries = 5

type DocumentBase struct {
	Key key.Key

	clear   func()
	data    any
	version noptions.Version

	DocumentStore diface.ICollection
	cache         diface.ICache
}

// Init performs an in-place initialization of a DocumentBase.
func (d *DocumentBase) Init(
	data any,
	clear func(),
	store diface.ICollection,
	key key.Key,
) {
	d.clear = clear
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
	cache diface.ICache,
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

// Create  data and version in the database.
func (d *DocumentBase) Create() error {
	version, err := d.DocumentStore.Set(
		d.Key,
		noptions.WithSource(d.data),
	)
	if err != nil {
		return err
	}
	d.version = version
	return nil
}

// VersionCache is a cache of a version and its data structure.
type VersionCache struct {
	Version any
	Data    any
}

// Load loads data and version from the database.
func (d *DocumentBase) Load() error {
	d.clear()
	if ok := d.cache.GetCache(d.Key, &VersionCache{
		Version: &d.version,
		Data:    d.data,
	}); !ok {
		if version, err := d.DocumentStore.Get(
			d.Key,
			noptions.WithDestination(d.data),
		); err != nil {
			return err
		} else {
			d.version = version
			d.cache.SetCache(d.Key, &VersionCache{
				Version: d.version,
				Data:    d.data,
			})
		}
	}
	return nil
}

// Save saves data to the database.
// compare the version in the database and swap it.
func (d *DocumentBase) Save() error {
	version, err := d.DocumentStore.Set(
		d.Key,
		noptions.WithSource(d.data),
		noptions.WithVersion(d.version),
	)
	if err != nil {
		return err
	}
	d.version = version
	d.cache.DeleteCache(d.Key)
	return nil
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

// Update change the data with the given function and CAS(compare and swap) save it to the database.
// If the function returns false, the update will be aborted.
// If the update CAS fails, the function will be retried up to MaxRetries times with a randomized backoff.
func (d *DocumentBase) Update(f func() bool) error {
	if err := d.doUpdate(f, func() error {
		return d.Save()
	}); err != nil {
		return err
	} else {
		return nil
	}
}

// Delete delete data from the database.
func (d *DocumentBase) Delete() error {
	if err := d.DocumentStore.Delete(d.Key); err != nil {
		return err
	}
	d.cache.DeleteCache(d.Key)
	return nil
}
