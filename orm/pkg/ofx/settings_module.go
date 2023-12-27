package ofx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

type SettingsParams struct {
	fx.In

	DatabaseURL string `name:"DatabaseURL"`
	CacheURL    string `name:"CacheURL"`
}

type SettingsResult struct {
	fx.Out

	DocumentURL string `name:"DatabaseURL" envconfig:"DATABASE_URL" default:"mongodb://localhost:27017"`
	CacheURL    string `name:"CacheURL" envconfig:"CACHE_URL" default:"redis://localhost:6379" `
}

func (sr *SettingsResult) LoadFromEnv() (err error) {
	err = utility.Load(sr)
	return
}

var SettingsModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
