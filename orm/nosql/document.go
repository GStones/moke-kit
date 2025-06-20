package nosql

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

const (
	// MaxRetries is the maximum number of retries for update operations
	MaxRetries = 5
	// DefaultCacheTTL is the default cache TTL for read-through caching
	DefaultCacheTTL = 30 * time.Minute
)

// DocumentBase represents a base document structure for NoSQL operations
type DocumentBase struct {
	Key key.Key

	clear   func()
	data    any
	version noptions.Version

	DocumentStore diface.ICollection
	cache         diface.ICache
	ctx           context.Context
}

// Init performs an in-place initialization of a DocumentBase.
func (d *DocumentBase) Init(
	ctx context.Context,
	data any,
	clear func(),
	store diface.ICollection,
	key key.Key,
) {
	defaultCache := diface.DefaultDocumentCache()
	d.InitWithCache(ctx, data, clear, store, key, defaultCache)
}

// InitWithCache performs an in-place initialization of a DocumentBase with cache.
func (d *DocumentBase) InitWithCache(
	ctx context.Context,
	data any,
	clear func(),
	store diface.ICollection,
	key key.Key,
	cache diface.ICache,
) {
	d.ctx = ctx
	d.clear = clear
	d.data = data
	d.DocumentStore = store
	d.Key = key
	d.cache = cache
}

// Clear clears all data on this DocumentBase.
func (d *DocumentBase) Clear() {
	d.version = noptions.NoVersion
	d.clear()
}

// Create  data and version in the database.
func (d *DocumentBase) Create() error {
	version, err := d.DocumentStore.Set(
		d.ctx,
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

// Load implements Read-Through caching
func (d *DocumentBase) Load() error {
	d.clear()
	cache := &VersionCache{
		Version: &d.version,
		Data:    d.data,
	}

	// Try cache first
	if d.cache.GetCache(d.ctx, d.Key, cache) {
		return nil
	}

	// Cache miss - load from database
	version, err := d.DocumentStore.Get(
		d.ctx,
		d.Key,
		noptions.WithDestination(d.data),
	)
	if err != nil {
		return err
	}

	d.version = version
	// Update cache after loading from database
	d.cache.SetCache(d.ctx, d.Key, &VersionCache{
		Version: d.version,
		Data:    d.data,
	}, DefaultCacheTTL)

	return nil
}

// Save implements synchronous write with cache update
func (d *DocumentBase) Save() error {
	// 直接同步写入数据库
	version, err := d.DocumentStore.Set(
		d.ctx,
		d.Key,
		noptions.WithSource(d.data),
		noptions.WithVersion(d.version),
	)
	if err != nil {
		return err
	}
	d.version = version

	// 更新缓存
	d.cache.SetCache(d.ctx, d.Key, &VersionCache{
		Version: d.version,
		Data:    d.data,
	}, DefaultCacheTTL)

	return nil
}

func (d *DocumentBase) doUpdate(f func() bool, u func() error) error {
	var lastErr error
	for r := 0; r < MaxRetries; r++ {
		if !f() {
			return nerrors.ErrUpdateLogicFailed
		}

		if err := u(); err == nil {
			return nil
		} else {
			lastErr = err
			// Exponential backoff with jitter
			backoff := time.Duration(math.Pow(2, float64(r))) * time.Millisecond
			jitter := time.Duration(rand.Float64() * float64(backoff))
			time.Sleep(backoff + jitter)

			if err := d.Load(); err != nil {
				return err
			}
		}
	}
	if lastErr != nil {
		return errors.Wrap(nerrors.ErrTooManyRetries, lastErr.Error())
	}
	return errors.Wrap(nerrors.ErrTooManyRetries, "no underlying error")
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
	if err := d.DocumentStore.Delete(d.ctx, d.Key); err != nil {
		return err
	}
	d.cache.DeleteCache(d.ctx, d.Key)
	return nil
}
