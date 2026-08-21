package diface

import (
	"context"
	"time"

	"github.com/gstones/moke-kit/orm/nosql/key"
)

// ICache provides a HASH-oriented cache for Document objects.
// Values are stored as Redis HASH field maps (or an equivalent in-memory map).
type ICache interface {
	// GetCache retrieves cache fields. When fields is empty, all fields are returned.
	// A cache miss returns nil or an empty map.
	GetCache(ctx context.Context, key key.Key, fields ...string) map[string]any
	// SetCache writes cache fields and optionally sets an expiration.
	SetCache(ctx context.Context, key key.Key, data map[string]any, expire time.Duration) error
	// DeleteCache deletes a Document from the cache.
	DeleteCache(ctx context.Context, key key.Key)
}

type defaultDocumentCache struct{}

var _ ICache = (*defaultDocumentCache)(nil)

// DefaultDocumentCache returns a new no-op ICache.
func DefaultDocumentCache() ICache {
	return &defaultDocumentCache{}
}

func (c *defaultDocumentCache) GetCache(ctx context.Context, key key.Key, fields ...string) map[string]any {
	return nil
}

func (c *defaultDocumentCache) SetCache(ctx context.Context, key key.Key, data map[string]any, expire time.Duration) error {
	return nil
}

func (c *defaultDocumentCache) DeleteCache(ctx context.Context, key key.Key) {
}
