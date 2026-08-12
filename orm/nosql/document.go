package nosql

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

const (
	// MaxRetries is the maximum number of retries for update operations.
	MaxRetries = 5
	// DefaultCacheTTL is the default cache TTL for read-through caching.
	DefaultCacheTTL = 30 * time.Minute
	// DefaultWriteBackDelay is the default async write-back delay.
	DefaultWriteBackDelay = 500 * time.Millisecond
	// ExpireRangeMin is the minimum jittered cache expiration.
	ExpireRangeMin = 6 * time.Hour
	// ExpireRangeMax is the maximum jittered cache expiration.
	ExpireRangeMax = 12 * time.Hour
	// WriteBackTopic is the MQ topic used for delayed Mongo write-back.
	WriteBackTopic = "nats://writeback"
)

// WriteBackPayload is the delayed write-back message payload.
type WriteBackPayload struct {
	CollectionName string           `json:"collection"`
	Key            string           `json:"key"`
	Data           json.RawMessage  `json:"data"`
	Version        noptions.Version `json:"version"`
}

// WriteBackOptions configures delayed write-back behavior.
type WriteBackOptions struct {
	Enabled bool
	Delay   time.Duration
	MQ      miface.MessageQueue
}

// DefaultWriteBackOptions returns disabled write-back options.
func DefaultWriteBackOptions() WriteBackOptions {
	return WriteBackOptions{
		Enabled: false,
		Delay:   DefaultWriteBackDelay,
		MQ:      nil,
	}
}

// Validate validates write-back options.
func (opts WriteBackOptions) Validate() error {
	if opts.Enabled && opts.MQ == nil {
		return errors.New("MQ client is required when WriteBack is enabled")
	}
	if opts.Delay < 0 {
		return errors.New("WriteBack delay cannot be negative")
	}
	return nil
}

// DocumentBase represents a base document structure for NoSQL operations.
type DocumentBase struct {
	Key key.Key

	clear   func()
	data    any
	version noptions.Version

	DocumentStore diface.ICollection
	cache         diface.ICache
	ctx           context.Context

	writeBack WriteBackOptions
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
	d.writeBack = DefaultWriteBackOptions()
}

// EnableWriteBackWithMQ enables delayed MQ write-back.
func (d *DocumentBase) EnableWriteBackWithMQ(mqClient miface.MessageQueue, delay time.Duration) error {
	if mqClient == nil {
		return errors.New("MQ client cannot be nil")
	}
	d.writeBack.Enabled = true
	d.writeBack.MQ = mqClient
	if delay > 0 {
		d.writeBack.Delay = delay
	}
	return d.writeBack.Validate()
}

// DisableWriteBack disables delayed write-back.
func (d *DocumentBase) DisableWriteBack() {
	d.writeBack.Enabled = false
	d.writeBack.MQ = nil
}

// Clear clears all data on this DocumentBase.
func (d *DocumentBase) Clear() {
	d.version = noptions.NoVersion
	d.clear()
}

func randomExpiration() time.Duration {
	return ExpireRangeMin + time.Duration(rand.Int63n(int64(ExpireRangeMax-ExpireRangeMin)))
}

func (d *DocumentBase) cacheEnvelope() (map[string]any, error) {
	raw, err := json.Marshal(d.data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal document for cache")
	}
	return map[string]any{
		cacheFieldVersion: d.version,
		cacheFieldData:    string(raw),
	}, nil
}

func (d *DocumentBase) updateCache() error {
	envelope, err := d.cacheEnvelope()
	if err != nil {
		return err
	}
	return d.cache.SetCache(d.ctx, d.Key, envelope, randomExpiration())
}

func (d *DocumentBase) loadFromCache(m map[string]any) error {
	ver, ok := parseCacheVersion(m[cacheFieldVersion])
	if !ok {
		return errors.New("cache entry missing version")
	}
	raw, err := asJSONBytes(m[cacheFieldData])
	if err != nil {
		return errors.Wrap(err, "cache entry missing data")
	}
	if err := json.Unmarshal(raw, d.data); err != nil {
		return errors.Wrap(err, "failed to unmarshal cached document")
	}
	d.version = ver
	return nil
}

func (d *DocumentBase) scheduleWriteBack(raw json.RawMessage, version noptions.Version) error {
	payload := &WriteBackPayload{
		CollectionName: d.DocumentStore.GetName(),
		Key:            d.Key.String(),
		Data:           raw,
		Version:        version,
	}
	return d.writeBack.MQ.Publish(WriteBackTopic, miface.WithJSON(payload))
}

// Create data and version in the database.
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
	return d.updateCache()
}

// Load implements Read-Through caching through Redis HASH fields.
func (d *DocumentBase) Load() error {
	d.clear()

	if cached := d.cache.GetCache(d.ctx, d.Key); len(cached) > 0 {
		if err := d.loadFromCache(cached); err == nil {
			return nil
		}
		// Corrupt/incomplete cache entries fall through to the database.
		d.cache.DeleteCache(d.ctx, d.Key)
		d.clear()
	}

	version, err := d.DocumentStore.Get(
		d.ctx,
		d.Key,
		noptions.WithDestination(d.data),
	)
	if err != nil {
		return err
	}
	d.version = version
	return d.updateCache()
}

// Save implements synchronous write with HASH cache update.
func (d *DocumentBase) Save() error {
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
	return d.updateCache()
}

// SaveAsync optimistically updates the HASH cache and schedules a delayed DB write-back.
// The DB CAS uses the previous version; the cache stores previous+1.
func (d *DocumentBase) SaveAsync() error {
	if d.version == noptions.NoVersion {
		return errors.New("cannot async-save a document without a version")
	}

	raw, err := json.Marshal(d.data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal document for async save")
	}

	prev := d.version
	d.version = prev + 1
	if err := d.updateCache(); err != nil {
		d.version = prev
		return err
	}

	if !d.writeBack.Enabled || d.writeBack.MQ == nil {
		return nil
	}

	delay := d.writeBack.Delay
	go func() {
		if delay > 0 {
			time.Sleep(delay)
		}
		_ = d.scheduleWriteBack(raw, prev)
	}()
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
			backoff := time.Duration(math.Pow(2, float64(r))) * time.Millisecond
			jitter := time.Duration(rand.Float64() * float64(backoff))
			time.Sleep(backoff + jitter)

			if err := d.Load(); err != nil {
				return errors.Wrap(err, "failed to reload during update retry")
			}
		}
	}
	if lastErr != nil {
		return errors.Wrap(nerrors.ErrTooManyRetries, lastErr.Error())
	}
	return errors.Wrap(nerrors.ErrTooManyRetries, "no underlying error")
}

// Update changes data with CAS save and retries.
func (d *DocumentBase) Update(f func() bool) error {
	return d.doUpdate(f, func() error {
		return d.Save()
	})
}

// UpdateAsync mutates in-memory data, updates HASH cache, and schedules write-back.
func (d *DocumentBase) UpdateAsync(f func() bool) error {
	if !f() {
		return nerrors.ErrUpdateLogicFailed
	}
	return d.SaveAsync()
}

// Delete deletes data from the database and cache.
func (d *DocumentBase) Delete() error {
	if err := d.DocumentStore.Delete(d.ctx, d.Key); err != nil {
		return err
	}
	d.cache.DeleteCache(d.ctx, d.Key)
	return nil
}
