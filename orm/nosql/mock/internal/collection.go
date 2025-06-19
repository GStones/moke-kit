package internal

import (
	"context"
	"encoding/json"

	"github.com/stretchr/testify/mock"

	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

type MockCollection struct {
	mock.Mock
	data    []byte
	version noptions.Version
}

func NewMockCollection(name string) *MockCollection {
	mc := &MockCollection{
		version: 0,
	}
	mc.On("GetName").Return(name)
	mc.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, key key.Key, opts ...noptions.Option) noptions.Version {
		mc.version++
		return mc.version
	}, nil)
	mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, key key.Key, opts ...noptions.Option) noptions.Version {
		if mc.version == 0 {
			mc.version = 1 // 确保至少返回版本1
		}
		return mc.version
	}, nil)
	mc.On("Delete", mock.Anything, mock.Anything).Return(nil)
	mc.On("Incr", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)
	return mc
}

func (m *MockCollection) GetName() string {
	ret := m.Called()
	r0 := ret.Get(0).(string)
	return r0
}

func (m *MockCollection) Set(ctx context.Context, key key.Key, opts ...noptions.Option) (noptions.Version, error) {
	options, err := noptions.NewOptions(opts...)
	if err != nil {
		return noptions.NoVersion, err
	}
	if options.Source != nil {
		jsonData, err := json.Marshal(options.Source)
		if err != nil {
			return noptions.NoVersion, err
		}
		m.data = jsonData
	}
	
	// 直接增加版本号
	m.version++
	
	ret := m.Called(ctx, key, opts)
	return m.version, ret.Error(1)
}

func (m *MockCollection) Get(
	ctx context.Context,
	key key.Key,
	opts ...noptions.Option,
) (noptions.Version, error) {
	options, err := noptions.NewOptions(opts...)
	if err != nil {
		return noptions.NoVersion, err
	}

	if options.Destination != nil && m.data != nil {
		json.Unmarshal(m.data, options.Destination)
	}

	ret := m.Called(ctx, key, opts)
	
	// 确保返回有效的版本号
	if m.version == 0 {
		m.version = 1
	}
	
	return m.version, ret.Error(1)
}

func (m *MockCollection) Delete(ctx context.Context, key key.Key) error {
	ret := m.Called(ctx, key)
	r0 := ret.Error(0)
	return r0
}

func (m *MockCollection) Incr(ctx context.Context, key key.Key, field string, amount int32) (int64, error) {
	ret := m.Called(ctx, key, field, amount)
	r0 := ret.Get(0).(int64)
	r1 := ret.Error(1)
	return r0, r1
}
