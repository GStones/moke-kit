package diface

import (
	"moke-kit/gorm/nosql/key"
	"moke-kit/gorm/nosql/noptions"
)

// IDocumentProvider knows how to open nosql stores by name.
type IDocumentProvider interface {
	OpenDbDriver(name string) (ICollection, error)
	Shutdown() error
}

// ICollection is an abstract container of NoSQL documents located by keys.
type ICollection interface {
	// GetName Name returns the name of this ICollection.
	GetName() string

	// Set creates or overwrites a nosql with the given keys and returns its version.  Use WithVersion to ensure this
	// function is updating the version of the nosql that you expect.  If you don't use WithVersion then this
	// function expects there to be no nosql.  If you want to set the nosql no matter what then use
	// WithAnyVersion.
	Set(key key.Key, opts ...noptions.Option) (noptions.Version, error)

	// Get loads an existing nosql from the nosql store and returns its cas.  If no such nosql exists then
	// this function fails.  Use WithTTL to update the nosql's expiration time.
	Get(key key.Key, opts ...noptions.Option) (noptions.Version, error)
}
