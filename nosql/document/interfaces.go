package document

import (
	"math"

	"github.com/gstones/platform/services/common/dupblock"
)

type Version = uint64

const (
	NoVersion Version = math.MaxUint64
)

// StartingDocument describes a single document to be loaded into a DocumentStoreProvider.  Generally used in mock
// implementations but may be useful if / when we want to initialize document stores on service startup.
type StartingDocument struct {
	// Store contains the name of the DocumentStore that will contain this document.
	Store string

	// Key contains the string of the keys for this document.
	Key string

	// Data to store in this document.  If a string / []byte is provided then the data will be taken as given.
	// Otherwise it will be serialized into JSON before being stored.
	Data interface{}
}

// DocumentStoreProvider knows how to open document stores by name.
type DocumentStoreProvider interface {
	AddStartingDocuments(documents []StartingDocument) error
	OpenDocumentStore(name string) (DocumentStore, error)
	Shutdown() error
}

// DocumentStore is an abstract container of NoSQL documents located by keys.
type DocumentStore interface {
	// Name returns the name of this DocumentStore.
	Name() string

	// Contains checks to see if a document with the given keys exists.
	Contains(key Key) (bool, error)

	// Return a list of all keys matching a given prefix. WithTimeout is supported for appropriate backend(s).
	// Additionally, supports offset and limit if desired for some reason.
	ListKeys(prefix string, opts ...ScanOption) ([]Key, error)

	// Set creates or overwrites a document with the given keys and returns its version.  Use WithVersion to ensure this
	// function is updating the version of the document that you expect.  If you don't use WithVersion then this
	// function expects there to be no document.  If you want to set the document no matter what then use
	// WithAnyVersion.
	Set(key Key, opts ...Option) (Version, error)

	// Get loads an existing document from the document store and returns its cas.  If no such document exists then
	// this function fails.  Use WithTTL to update the document's expiration time.
	Get(key Key, opts ...Option) (Version, error)

	// Scan a filtered set of keys for one or more objects that match a specific query. WithDestination must be used in
	// the second parameter to provide a target pointer - which may either be an individual object or an array/slice.
	// Use MatchKeyValue to set the keys/value pair to be searched for and WithLimit (or WithNoLimit) to explicitly
	// provide a number of results to stop at. Current implementation can only handle matching against a single K/V
	// pair, but future extensions to the query variety are anticipated. WithTimeout is supported for appropriate
	// backend(s).
	Scan(prefix string, dest Option, scanOpts ...ScanOption) (int, error)

	// Remove removes an existing document from the document store.  Use WithVersion to ensure this function is
	// removing the version of the document that you expect.  Use WithAnyVersion to remove the document no matter what.
	Remove(key Key, opts ...Option) error

	// ApplyDUPBlock applies a document update block to the DocumentStore.
	ApplyDUPBlock(reader dupblock.Reader) error

	// SetField sets a field (denoted by path) on a document and returns its cas.  If no such document exists then this
	// function fails.  Use WithVersion to ensure this function is updating the version of the document that you expect.
	SetField(key Key, path string, opts ...Option) (Version, error)

	// AddWatcher adds a DocumentWatcher to the given keys.  Whenever changes to the document occur the watcher is told.
	AddWatcher(key Key, watcher DocumentWatcher)

	// RemoveWatcher removes a DocumentWatcher from the given keys.  It's the cleanup pair for AddWatcher.
	RemoveWatcher(key Key, watcher DocumentWatcher)

	// WatchersForKey returns a new slice of all DocumentWatcher instances for a keys.
	WatchersForKey(key Key) []DocumentWatcher

	// PushBack adds an element to the end of the array at the provided keys.
	PushBack(key Key, path string, item interface{}) error

	// Incr increments a numeric field by the provided amount.
	Incr(key Key, path string, amount int32) error

	// Touch touches a document, specifying a new expiry time for it.
	Touch(key Key, opts ...Option) (Version, error)
}

// DocumentWatcher receives notifications when a document is modified.
type DocumentWatcher interface {
	OnDocumentChanged(key Key, value interface{})
	OnDocumentExpired(key Key)
	OnDocumentDeleted(key Key)
}
