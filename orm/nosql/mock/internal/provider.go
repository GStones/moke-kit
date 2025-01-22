package internal

import (
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/diface"
)

type MockDriverProvider struct {
	logger *zap.Logger
}

func (dp *MockDriverProvider) Shutdown() error {
	return nil
}

func (dp *MockDriverProvider) OpenDbDriver(name string) (diface.ICollection, error) {
	mc := NewMockCollection(name)
	return mc, nil
}

func NewMockDriverProvider(
	logger *zap.Logger,
) *MockDriverProvider {
	m := &MockDriverProvider{logger: logger}
	return m
}
