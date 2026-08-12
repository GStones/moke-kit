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
	CollectionName string          `json:"collection"`
	Key            string          `json:"key"`
	Data           json.RawMessage `json:"data"`
	// Version is the optimistic target version after SaveAsync (not the CAS base).
	Version noptions.Version `json:"version"`
	// Epoch fences delete/recreate generations when workers can read cache.
	Epoch int64 `json:"epoch"`
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
	epoch   int64

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
	d.version = noptions.NoVersion
	d.epoch = 0
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
	d.epoch = 0
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
	envelope := map[string]any{
		cacheFieldVersion: d.version,
		cacheFieldData:    string(raw),
	}
	// Omit zero epochs so Load-from-DB after restart does not fence in-flight write-backs.
	if d.epoch != 0 {
		envelope[cacheFieldEpoch] = d.epoch
	}
	return envelope, nil
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
	if ep, ok := parseCacheVersion(m[cacheFieldEpoch]); ok {
		d.epoch = ep
	} else {
		d.epoch = 0
	}
	return nil
}

func scheduleWriteBack(
	mq miface.MessageQueue,
	collectionName string,
	keyStr string,
	raw json.RawMessage,
	version noptions.Version,
	epoch int64,
) error {
	if mq == nil {
		return errors.New("MQ client is nil")
	}
	payload := &WriteBackPayload{
		CollectionName: collectionName,
		Key:            keyStr,
		Data:           raw,
		Version:        version,
		Epoch:          epoch,
	}
	return mq.Publish(WriteBackTopic, miface.WithJSON(payload))
}

var errWriteBackStale = errors.New("write-back snapshot is stale")

// applyWriteBackSnapshot installs src and advances the DB CAS version up to targetVersion.
// Store Set only $inc version by 1, so a gapped target (e.g. 5 while DB is 1) is
// fast-forwarded with the same payload until DB version == targetVersion.
// Older/equal targets are dropped so out-of-order delayed messages cannot clobber newer data.
func applyWriteBackSnapshot(
	ctx context.Context,
	coll diface.ICollection,
	docKey key.Key,
	src any,
	targetVersion noptions.Version,
) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		var discard any
		current, err := coll.Get(ctx, docKey, noptions.WithDestination(&discard))
		if err != nil {
			return err
		}
		if targetVersion <= current {
			return errWriteBackStale
		}
		newVer, err := coll.Set(
			ctx,
			docKey,
			noptions.WithSource(src),
			noptions.WithVersion(current),
		)
		if err != nil {
			return err
		}
		if newVer >= targetVersion {
			return nil
		}
		// Version advanced but still behind the optimistic target; keep fast-forwarding.
	}
}

// Create data and version in the database.
// Each Create assigns a unique epoch so delayed write-backs from a prior
// delete/recreate generation can be fenced across DocumentBase instances.
func (d *DocumentBase) Create() error {
	prevEpoch := d.epoch
	d.epoch = time.Now().UnixNano()
	if d.epoch == prevEpoch {
		d.epoch++
	}
	version, err := d.DocumentStore.Set(
		d.ctx,
		d.Key,
		noptions.WithSource(d.data),
	)
	if err != nil {
		d.epoch = prevEpoch
		return err
	}
	d.version = version
	if err := d.updateCache(); err != nil {
		d.epoch = prevEpoch
		return err
	}
	return nil
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
	// DB has no persisted epoch. Preserve any existing cache fence so a concurrent
	// Create/SaveAsync generation is not wiped by this read-through rewrite.
	d.epoch = 0
	if cached := d.cache.GetCache(d.ctx, d.Key, cacheFieldEpoch); len(cached) > 0 {
		if ep, ok := parseCacheVersion(cached[cacheFieldEpoch]); ok {
			d.epoch = ep
		}
	}
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
// The write-back payload carries the optimistic target version (previous+1).
// When write-back is disabled, SaveAsync falls back to synchronous Save().
func (d *DocumentBase) SaveAsync() error {
	if d.version == noptions.NoVersion {
		return errors.New("cannot async-save a document without a version")
	}

	if !d.writeBack.Enabled || d.writeBack.MQ == nil {
		return d.Save()
	}

	raw, err := json.Marshal(d.data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal document for async save")
	}

	prev := d.version
	next := prev + 1
	d.version = next
	if err := d.updateCache(); err != nil {
		d.version = prev
		return err
	}

	// Capture publish dependencies so DisableWriteBack cannot nil-deref the goroutine.
	mq := d.writeBack.MQ
	store := d.DocumentStore
	cache := d.cache
	docKey := d.Key
	collectionName := store.GetName()
	keyStr := docKey.String()
	delay := d.writeBack.Delay
	epoch := d.epoch
	go func() {
		if delay > 0 {
			time.Sleep(delay)
		}
		// Version is the optimistic target version (next), not the CAS base.
		if err := scheduleWriteBack(mq, collectionName, keyStr, raw, next, epoch); err != nil {
			// Publish failed: best-effort sync fallback so data is not stuck in cache only.
			// Use a fresh timeout context; the request ctx is usually already cancelled.
			var src any
			if json.Unmarshal(raw, &src) != nil {
				return
			}
			fbCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if cacheEpochMismatch(fbCtx, cache, docKey, epoch) {
				return
			}
			_ = applyWriteBackSnapshot(fbCtx, store, docKey, src, next)
		}
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
