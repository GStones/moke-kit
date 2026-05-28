package diface

import (
	"context"
	"time"

	"github.com/gstones/moke-kit/orm/nosql/key"
)

// ICache provides a cache for Document objects.
type ICache interface {
	// GetCache Get retrieves a Document from the cache.
	GetCache(ctx context.Context, key key.Key, fields ...string) map[string]any
	// SetCache Set sets a Document in the cache.
	SetCache(ctx context.Context, key key.Key, data map[string]any, expire time.Duration) error
	// DeleteCache Delete deletes a Document from the cache.
	DeleteCache(ctx context.Context, key key.Key)
}
type defaultDocumentCache struct {
}

var _ ICache = (*defaultDocumentCache)(nil)

// DefaultDocumentCache returns a new ICache.
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
