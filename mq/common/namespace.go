package common

import (
	"strings"

	"github.com/gstones/moke-kit/utility"
)

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
