package diface

import (
	"moke-kit/nosql/document/key"
)

type Version = int64

const (
	NoVersion Version = 0
)

// IDocumentProvider knows how to open document stores by name.
type IDocumentProvider interface {
	OpenDbDriver(name string) (ICollection, error)
	Shutdown() error
}

// ICollection is an abstract container of NoSQL documents located by keys.
type ICollection interface {
	// GetName Name returns the name of this ICollection.
	GetName() string

	// Set creates or overwrites a document with the given keys and returns its version.  Use WithVersion to ensure this
	// function is updating the version of the document that you expect.  If you don't use WithVersion then this
	// function expects there to be no document.  If you want to set the document no matter what then use
	// WithAnyVersion.
	Set(key key.Key, opts ...Option) (Version, error)

	// Get loads an existing document from the document store and returns its cas.  If no such document exists then
	// this function fails.  Use WithTTL to update the document's expiration time.
	Get(key key.Key, opts ...Option) (Version, error)
}
