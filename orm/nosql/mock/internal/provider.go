package internal

import (
	"sync"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/diface"
)

type MockDriverProvider struct {
	logger *zap.Logger
	mu     sync.Mutex
	cols   map[string]*MockCollection
}

func (dp *MockDriverProvider) Shutdown() error {
	return nil
}

func (dp *MockDriverProvider) OpenDbDriver(name string) (diface.ICollection, error) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	if dp.cols == nil {
		dp.cols = make(map[string]*MockCollection)
	}
	if mc, ok := dp.cols[name]; ok {
		return mc, nil
	}
	mc := NewMockCollection(name)
	dp.cols[name] = mc
	return mc, nil
}

func NewMockDriverProvider(
	logger *zap.Logger,
) *MockDriverProvider {
	return &MockDriverProvider{
		logger: logger,
		cols:   make(map[string]*MockCollection),
	}
}
