package utility

import "go.uber.org/atomic"

// namespace is process-wide multi-tenant prefix state (one deployment per process).
var namespace = atomic.NewString("")

// Namespace returns the current global deployment namespace.
func Namespace() string {
	return namespace.Load()
}

// SetNamespace sets the global deployment namespace used by ORM keys and MQ topics.
func SetNamespace(ns string) {
	namespace.Store(ns)
}
