package siface

import (
	"time"

	"agones.dev/agones/pkg/sdk"
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
}
