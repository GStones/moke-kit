package siface

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
	Alpha() *agone.Alpha
}
