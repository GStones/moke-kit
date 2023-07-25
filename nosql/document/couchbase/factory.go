package couchbase

import (
	"github.com/gstones/platform/services/common/nosql/document"
	"github.com/gstones/platform/services/common/nosql/document/couchbase/internal"
	"go.uber.org/zap"
)

func NewDocumentStoreProvider(config ClusterConfig, l *zap.Logger) (document.DocumentStoreProvider, error) {
	return internal.NewDocumentStoreProvider(config.ConnUrl, config.Username, config.Password, l)
}
