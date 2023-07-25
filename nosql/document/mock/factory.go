package mock

import (
	"github.com/gstones/platform/services/common/nosql/document"
	"github.com/gstones/platform/services/common/nosql/document/mock/internal"
)

type DocumentStoreProvider = internal.DocumentStoreProvider

func NewDocumentStoreProvider() (document.DocumentStoreProvider, error) {
	return internal.NewDocumentStoreProvider()
}
