package aiface

import (
	"time"

	"agones.dev/agones/pkg/sdk"
	agone "agones.dev/agones/sdks/go"
)

type IAgones interface {
	Init() error
	Ready() error
	Health() error
	Allocate() error
	Shutdown() error
	Reserve(d time.Duration) error
	SetLabel(key, value string) error
	SetAnnotation(key, value string) error
	GameServer() (*sdk.GameServer, error)
	WatchGameServer(_ agone.GameServerCallback) error
	CounterList() ICounterList
}

type ICounterList interface {
	GetCounterCount(key string) (int64, error)
	IncrementCounter(key string, amount int64) error
	DecrementCounter(key string, amount int64) error
	SetCounterCount(key string, amount int64) error
	GetCounterCapacity(key string) (int64, error)
	SetCounterCapacity(key string, amount int64) error
	GetListCapacity(key string) (int64, error)
	SetListCapacity(key string, amount int64) error
	ListContains(key, value string) (bool, error)
	GetListLength(key string) (int, error)
	GetListValues(key string) ([]string, error)
	AppendListValue(key, value string) error
	DeleteListValue(key, value string) error
}
