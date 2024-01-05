package agones

import (
	"time"

	"agones.dev/agones/pkg/sdk"
)

type Default struct {
}

func (d *Default) Init() error {
	return nil
}

func (d *Default) Ready() error {
	return nil
}

func (d *Default) Health() error {
	return nil
}

func (d *Default) Allocate() error {
	return nil
}

func (d *Default) Shutdown() error {
	return nil
}

func (d *Default) Reserve(_ time.Duration) error {
	return nil
}

func (d *Default) SetLabel(_, _ string) error {
	return nil
}

func (d *Default) SetAnnotation(_, _ string) error {
	return nil
}

func (d *Default) GameServer() (*sdk.GameServer, error) {
	return nil, nil
}

func NewDefault() *Default {
	return &Default{}
}
