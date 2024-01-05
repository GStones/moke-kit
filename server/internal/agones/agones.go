package agones

import (
	"time"

	"agones.dev/agones/pkg/sdk"
	agone "agones.dev/agones/sdks/go"
)

// TickDuration is the duration between health checks
const TickDuration = 2 * time.Second

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
	go a.health()
	return nil
}

func (a *Agones) health() {
	tick := time.NewTicker(TickDuration)
	for {
		<-tick.C
		if err := a.sdk.Health(); err != nil {
			continue
		}
	}

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
