package common

import (
	"strings"

	"go.uber.org/atomic"
)

// This is a convenience mechanism for supporting multi-tenant usage of
// a single Kafka deployment. When referencing a topic, abstractions within this
// library will prefix the topic name with the namespace.
//
// It's important for consumers of this library to understand that this mechanism
// relies on global state. This is fine as long as the library's usage is
// confined to a single process.
var namespace = atomic.NewString("")

// Namespace returns the current global fxapp namespace.
func Namespace() string {
	return namespace.Load()
}

// SetNamespace sets the global fxapp namespace.
func SetNamespace(ns string) {
	namespace.Store(ns)
}

const (
	NamespaceSep = "."
)

func NamespaceTopic(topic string) string {
	namespace := Namespace()
	if namespace != "" {
		return strings.Join([]string{namespace, topic}, NamespaceSep)
	} else {
		return topic
	}
}
