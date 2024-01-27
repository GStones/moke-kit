package sdk

import (
	"time"

	"agones.dev/agones/pkg/sdk"
	agone "agones.dev/agones/sdks/go"

	"github.com/gstones/moke-kit/3rd/agones/aiface"
)

// Agones is a wrapper around the agones sdk
type Agones struct {
	sdk *agone.SDK
}

func (a *Agones) Init() error {
	s, err := agone.NewSDK()
	if err != nil {
		return err
	}
	a.sdk = s
	return nil
}

func (a *Agones) Ready() error {
	return a.sdk.Ready()
}

func (a *Agones) Health() error {
	return a.sdk.Health()
}

func (a *Agones) Allocate() error {
	return a.sdk.Allocate()
}

func (a *Agones) Shutdown() error {
	return a.sdk.Shutdown()
}

func (a *Agones) Reserve(d time.Duration) error {
	return a.sdk.Reserve(d)
}

func (a *Agones) SetLabel(key, value string) error {
	return a.sdk.SetLabel(key, value)
}

func (a *Agones) SetAnnotation(key, value string) error {
	return a.sdk.SetAnnotation(key, value)
}

func (a *Agones) GameServer() (*sdk.GameServer, error) {
	return a.sdk.GameServer()
}

func (a *Agones) WatchGameServer(cb agone.GameServerCallback) error {
	if err := a.sdk.WatchGameServer(cb); err != nil {
		return err
	}
	return nil
}

func (a *Agones) Alpha() aiface.IAlpha {
	return a.sdk.Alpha()
}
