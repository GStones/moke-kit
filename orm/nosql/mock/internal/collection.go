package internal

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/stretchr/testify/mock"

	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

type docEntry struct {
	data    []byte
	version noptions.Version
}

type MockCollection struct {
	mock.Mock
	mu   sync.Mutex
	docs map[string]*docEntry
}

func NewMockCollection(name string) *MockCollection {
	mc := &MockCollection{
		docs: make(map[string]*docEntry),
	}
	mc.On("GetName").Return(name)
	mc.On("Set", mock.Anything, mock.Anything, mock.Anything).Maybe()
	mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Maybe()
	mc.On("Delete", mock.Anything, mock.Anything).Maybe()
	mc.On("Incr", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	return mc
}

func (m *MockCollection) GetName() string {
	ret := m.Called()
	return ret.Get(0).(string)
}

func (m *MockCollection) Set(ctx context.Context, k key.Key, opts ...noptions.Option) (noptions.Version, error) {
	m.Called(ctx, k, opts)
	options, err := noptions.NewOptions(opts...)
	if err != nil {
		return noptions.NoVersion, err
	}
	if options.Source == nil {
		return noptions.NoVersion, nerrors.ErrSourceIsNil
	}
	jsonData, err := json.Marshal(options.Source)
	if err != nil {
		return noptions.NoVersion, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	id := k.String()
	entry, exists := m.docs[id]

	if !options.AnyVersion && options.Version == noptions.NoVersion {
		if exists {
			return noptions.NoVersion, nerrors.ErrVersionNotMatch
		}
		m.docs[id] = &docEntry{data: jsonData, version: 1}
		return 1, nil
	}
	if !options.AnyVersion {
		if !exists || entry.version != options.Version {
			return noptions.NoVersion, nerrors.ErrVersionNotMatch
		}
		entry.data = jsonData
		entry.version++
		return entry.version, nil
	}
	if !exists {
		m.docs[id] = &docEntry{data: jsonData, version: 1}
		return 1, nil
	}
	entry.data = jsonData
	entry.version++
	return entry.version, nil
}

func (m *MockCollection) Get(
	ctx context.Context,
	k key.Key,
	opts ...noptions.Option,
) (noptions.Version, error) {
	m.Called(ctx, k, opts)
	options, err := noptions.NewOptions(opts...)
	if err != nil {
		return noptions.NoVersion, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	entry, exists := m.docs[k.String()]
	if !exists {
		return noptions.NoVersion, nerrors.ErrNotFound
	}
	if options.Version != noptions.NoVersion && entry.version != options.Version {
		return noptions.NoVersion, nerrors.ErrNotFound
	}
	if options.Destination != nil && entry.data != nil {
		if err := json.Unmarshal(entry.data, options.Destination); err != nil {
			return noptions.NoVersion, err
		}
	}
	return entry.version, nil
}

func (m *MockCollection) Delete(ctx context.Context, k key.Key) error {
	m.Called(ctx, k)
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.docs[k.String()]; !exists {
		return nerrors.ErrNotFound
	}
	delete(m.docs, k.String())
	return nil
}

func (m *MockCollection) Incr(ctx context.Context, k key.Key, field string, amount int32) (int64, error) {
	m.Called(ctx, k, field, amount)
	m.mu.Lock()
	defer m.mu.Unlock()

	id := k.String()
	entry, exists := m.docs[id]
	var values map[string]int64
	if !exists {
		values = map[string]int64{field: int64(amount)}
		data, _ := json.Marshal(values)
		m.docs[id] = &docEntry{data: data, version: 1}
		return int64(amount), nil
	}
	values = map[string]int64{}
	_ = json.Unmarshal(entry.data, &values)
	values[field] += int64(amount)
	data, _ := json.Marshal(values)
	entry.data = data
	return values[field], nil
}
