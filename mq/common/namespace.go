package common

import (
	"strings"

	"github.com/gstones/moke-kit/utility"
)

// Namespace is a convenience mechanism for supporting multi-tenant usage of
// a single message queue deployment. When referencing a topic, abstractions
// within this library will prefix the topic name with the namespace.
//
// It's important for consumers of this library to understand that this mechanism
// relies on process-wide state shared with ORM keys via utility.Namespace.

const (
	NamespaceSep = "."
)

// Namespace returns the current global deployment namespace.
func Namespace() string {
	return utility.Namespace()
}

// SetNamespace sets the process-wide namespace used to prefix MQ topics.
func SetNamespace(ns string) {
	utility.SetNamespace(ns)
}

func NamespaceTopic(topic string) string {
	namespace := Namespace()
	if namespace != "" {
		return strings.Join([]string{namespace, topic}, NamespaceSep)
	}
	return topic
}
