package ofx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

type SettingsParams struct {
	fx.In

	DatabaseURL      string `name:"DatabaseURL"`
	DatabaseUser     string `name:"DatabaseUser"`     // will replace  database url username
	DatabasePassword string `name:"DatabasePassword"` // will replace  database url password
	CacheURL         string `name:"CacheURL"`
	CacheUser        string `name:"CacheUser"`     // will replace  cache url username
	CachePassword    string `name:"CachePassword"` // will replace  cache url password
}

type SettingsResult struct {
	fx.Out

	DocumentURL      string `name:"DatabaseURL" envconfig:"DATABASE_URL" default:"mongodb://localhost:27017"`
	DatabaseUser     string `name:"DatabaseUser" envconfig:"DATABASE_USER" default:""`
	DatabasePassword string `name:"DatabasePassword" envconfig:"DATABASE_PASSWORD" default:""`
	CacheURL         string `name:"CacheURL" envconfig:"CACHE_URL" default:"redis://localhost:6379"`
	CacheUser        string `name:"CacheUser" envconfig:"CACHE_USER" default:""`
	CachePassword    string `name:"CachePassword" envconfig:"CACHE_PASSWORD" default:""`
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
