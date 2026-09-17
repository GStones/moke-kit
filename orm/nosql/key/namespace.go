package key

import "github.com/gstones/moke-kit/utility"

// Namespace is a convenience mechanism for supporting multi-tenant usage of
// a single NoSQL deployment. When referencing a nosql, abstractions
// within this library will prefix the nosql's coordinates with the namespace.
//
// It's important for consumers of this library to understand that this mechanism
// relies on process-wide state shared with MQ topics via utility.Namespace.

// Namespace returns the current global deployment namespace.
func Namespace() string {
	return utility.Namespace()
}

// NamespaceKeyPrefix returns the nosql key prefix for the current namespace.
func NamespaceKeyPrefix() string {
	ns := utility.Namespace()
	if ns == "" {
		return ""
	}
	return KeySeparator + ns + KeySeparator
}

// SetNamespace sets the process-wide namespace used to prefix nosql keys.
func SetNamespace(ns string) {
	utility.SetNamespace(ns)
}
