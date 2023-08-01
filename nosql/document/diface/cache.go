package diface

import (
	"moke-kit/nosql/document/key"
)

// DocumentCache provides a cache for Document objects.
type DocumentCache interface {
	// Get retrieves a Document from the cache.
	GetCache(key key.Key, doc interface{}) bool
	// Set sets a Document in the cache.
	SetCache(key key.Key, doc interface{})
	// Delete deletes a Document from the cache.
	DeleteCache(key key.Key)
}
type defaultDocumentCache struct {
}

// NewDocumentCache returns a new DocumentCache.
func NewDocumentCache() DocumentCache {
	return &defaultDocumentCache{}
}

func (c *defaultDocumentCache) GetCache(key key.Key, doc interface{}) bool {
	return false
}

func (c *defaultDocumentCache) SetCache(key key.Key, doc interface{}) {
}

func (c *defaultDocumentCache) DeleteCache(key key.Key) {
}
