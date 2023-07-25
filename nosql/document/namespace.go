package document

// This is a convenience mechanism for supporting multi-tenant usage of
// a single NoSQL deployment. When referencing a document, abstractions
// within this library will prefix the document's coordinates with the namespace.
//
// It's important for consumers of this library to understand that this mechanism
// relies on global state. This is fine as long as the library's usage is
// confined to a single process.
var namespace = ""
var namespaceKeyPrefix = ""

// Namespace returns the current global application namespace.
func Namespace() string {
	return namespace
}

func NamespaceKeyPrefix() string {
	return namespaceKeyPrefix
}

// SetNamespace sets the global application namespace.
func SetNamespace(ns string) {
	namespace = ns
	namespaceKeyPrefix = KeySeparator + ns + KeySeparator
}
