package qfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

type SettingsParams struct {
	fx.In

	NatsUrl string `name:"NatsUrl"`
}

type SettingsResult struct {
	fx.Out

	NatsUrl string `name:"NatsUrl" envconfig:"NATS_URL" default:"nats://localhost:4222"`
}

func (ar *SettingsResult) LoadFromEnv() (err error) {
	err = utility.Load(ar)
	return
}

var SettingModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
