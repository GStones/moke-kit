package key

import "go.uber.org/atomic"

// This is a convenience mechanism for supporting multi-tenant usage of
// a single NoSQL deployment. When referencing a nosql, abstractions
// within this library will prefix the nosql's coordinates with the namespace.
//
// It's important for consumers of this library to understand that this mechanism
// relies on global state. This is fine as long as the library's usage is
// confined to a single process.
var namespace = atomic.NewString("")
var namespaceKeyPrefix = atomic.NewString("")

// Namespace returns the current global fxapp namespace.
func Namespace() string {
	return namespace.Load()
}

// NamespaceKeyPrefix returns the current global namespace key prefix.
func NamespaceKeyPrefix() string {
	return namespaceKeyPrefix.Load()
}

// SetNamespace sets the global fxapp namespace.
func SetNamespace(ns string) {
	namespace.Store(ns)
	namespaceKeyPrefix.Store(KeySeparator + ns + KeySeparator)
}
