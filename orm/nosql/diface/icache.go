package diface

import (
	"context"
	"time"

	"github.com/gstones/moke-kit/orm/nosql/key"
)

// ICache provides a cache for Document objects.
type ICache interface {
	// GetCache Get retrieves a Document from the cache.
	GetCache(ctx context.Context, key key.Key, doc any) bool
	// SetCache Set sets a Document in the cache.
	SetCache(ctx context.Context, key key.Key, doc any, expire time.Duration)
	// DeleteCache Delete deletes a Document from the cache.
	DeleteCache(ctx context.Context, key key.Key)
}
type defaultDocumentCache struct {
}

// DefaultDocumentCache returns a new ICache.
func DefaultDocumentCache() ICache {
	return &defaultDocumentCache{}
}

func (c *defaultDocumentCache) GetCache(ctx context.Context, key key.Key, doc any) bool {
	return false
}

func (c *defaultDocumentCache) SetCache(ctx context.Context, key key.Key, doc any, expire time.Duration) {
}

func (c *defaultDocumentCache) DeleteCache(ctx context.Context, key key.Key) {
}
