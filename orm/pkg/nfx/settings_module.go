package nfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility/uconfig"
)

type SettingsParams struct {
	fx.In

	DocumentURL   string `name:"DocumentURL"`
	RedisURL      string `name:"RedisURL"`
	RedisUser     string `name:"RedisUser"`
	RedisPassword string `name:"RedisPassword"`
}

type SettingsResult struct {
	fx.Out

	DocumentURL   string `name:"DocumentURL" envconfig:"DOCUMENT_URL" default:"mongodb://localhost:27017"`
	RedisURL      string `name:"RedisURL" envconfig:"REDIS_URL" default:"localhost:6379" `
	RedisUser     string `name:"RedisUser" envconfig:"REDIS_USER"`
	RedisPassword string `name:"RedisPassword" envconfig:"REDIS_PASSWORD"`
}

func (sr *SettingsResult) LoadFromEnv() (err error) {
	err = uconfig.Load(sr)
	return
}

var SettingsModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
