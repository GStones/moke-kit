package diface

import (
	"github.com/gstones/moke-kit/orm/nosql/key"
)

// ICache provides a cache for Document objects.
type ICache interface {
	// GetCache Get retrieves a Document from the cache.
	GetCache(key key.Key, doc any) bool
	// SetCache Set sets a Document in the cache.
	SetCache(key key.Key, doc any)
	// DeleteCache Delete deletes a Document from the cache.
	DeleteCache(key key.Key)
}
type defaultDocumentCache struct {
}

// DefaultDocumentCache returns a new ICache.
func DefaultDocumentCache() ICache {
	return &defaultDocumentCache{}
}

func (c *defaultDocumentCache) GetCache(key key.Key, doc any) bool {
	return false
}

func (c *defaultDocumentCache) SetCache(key key.Key, doc any) {
}

func (c *defaultDocumentCache) DeleteCache(key key.Key) {
}
