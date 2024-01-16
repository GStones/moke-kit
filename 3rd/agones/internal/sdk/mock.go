package sdk

import (
	"time"

	"agones.dev/agones/pkg/sdk"
	agone "agones.dev/agones/sdks/go"
)

type Mock struct {
}

func (d *Mock) Init() error {
	return nil
}

func (d *Mock) Ready() error {
	return nil
}

func (d *Mock) Health() error {
	return nil
}

func (d *Mock) Allocate() error {
	return nil
}

func (d *Mock) Shutdown() error {
	return nil
}

func (d *Mock) Reserve(_ time.Duration) error {
	return nil
}

func (d *Mock) SetLabel(_, _ string) error {
	return nil
}

func (d *Mock) SetAnnotation(_, _ string) error {
	return nil
}

func (d *Mock) GameServer() (*sdk.GameServer, error) {
	return nil, nil
}

func (d *Mock) WatchGameServer(_ agone.GameServerCallback) error {
	return nil
}

func (d *Mock) Alpha() *agone.Alpha {
	return nil
}
