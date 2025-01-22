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
	mc.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(mc.version+1, nil)
	mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(mc.version, nil)
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
	ret := m.Called(ctx, key, opts)
	r1 := ret.Error(1)
	if r1 == nil {
		m.version = m.version + 1
	}
	return m.version, r1
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

	if options.Destination != nil {
		json.Unmarshal(m.data, options.Destination)
	}

	ret := m.Called(ctx, key, opts)
	r0 := ret.Get(0).(noptions.Version)
	r1 := ret.Error(1)
	return r0, r1
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
