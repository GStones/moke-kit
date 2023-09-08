package diface

import (
	"github.com/gstones/moke-kit/orm/nosql/key"
)

// IDocumentCache provides a cache for Document objects.
type IDocumentCache interface {
	// GetCache Get retrieves a Document from the cache.
	GetCache(key key.Key, doc any) bool
	// SetCache Set sets a Document in the cache.
	SetCache(key key.Key, doc any)
	// DeleteCache Delete deletes a Document from the cache.
	DeleteCache(key key.Key)
}
type defaultDocumentCache struct {
}

// DefaultDocumentCache returns a new IDocumentCache.
func DefaultDocumentCache() IDocumentCache {
	return &defaultDocumentCache{}
}

func (c *defaultDocumentCache) GetCache(key key.Key, doc any) bool {
	return false
}

func (c *defaultDocumentCache) SetCache(key key.Key, doc any) {
}

func (c *defaultDocumentCache) DeleteCache(key key.Key) {
}
