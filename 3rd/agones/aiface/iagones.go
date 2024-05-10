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
	Alpha() IAlpha
}

type IAlpha interface {
	GetPlayerCapacity() (int64, error)
	SetPlayerCapacity(capacity int64) error
	PlayerConnect(id string) (bool, error)
	PlayerDisconnect(id string) (bool, error)
	GetPlayerCount() (int64, error)
	IsPlayerConnected(id string) (bool, error)
	GetConnectedPlayers() ([]string, error)
	GetCounterCount(key string) (int64, error)
	IncrementCounter(key string, amount int64) (bool, error)
	DecrementCounter(key string, amount int64) (bool, error)
	SetCounterCount(key string, amount int64) (bool, error)
	GetCounterCapacity(key string) (int64, error)
	SetCounterCapacity(key string, amount int64) (bool, error)
	GetListCapacity(key string) (int64, error)
	SetListCapacity(key string, amount int64) (bool, error)
	ListContains(key, value string) (bool, error)
	GetListLength(key string) (int, error)
	GetListValues(key string) ([]string, error)
	AppendListValue(key, value string) (bool, error)
	DeleteListValue(key, value string) (bool, error)
}
