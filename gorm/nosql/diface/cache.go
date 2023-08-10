package diface

import (
	"moke-kit/gorm/nosql/key"
)

// IDocumentCache provides a cache for Document objects.
type IDocumentCache interface {
	// GetCache Get retrieves a Document from the cache.
	GetCache(key key.Key, doc interface{}) bool
	// SetCache Set sets a Document in the cache.
	SetCache(key key.Key, doc interface{})
	// DeleteCache Delete deletes a Document from the cache.
	DeleteCache(key key.Key)
}
type defaultDocumentCache struct {
}

// DefaultDocumentCache returns a new IDocumentCache.
func DefaultDocumentCache() IDocumentCache {
	return &defaultDocumentCache{}
}

func (c *defaultDocumentCache) GetCache(key key.Key, doc interface{}) bool {
	return false
}

func (c *defaultDocumentCache) SetCache(key key.Key, doc interface{}) {
}

func (c *defaultDocumentCache) DeleteCache(key key.Key) {
}
