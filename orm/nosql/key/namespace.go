package key

import "github.com/gstones/moke-kit/utility"

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
