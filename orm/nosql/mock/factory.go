package mock

import (
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/mock/internal"
)

func NewMockDriverProvider(
	logger *zap.Logger,
) diface.IDocumentProvider {
	return internal.NewMockDriverProvider(logger)
}
