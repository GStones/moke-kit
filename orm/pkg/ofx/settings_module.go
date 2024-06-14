package ofx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

type SettingsParams struct {
	fx.In

	// DatabaseURL is the url of the database(mongodb).
	DatabaseURL string `name:"DatabaseURL"`
	// will replace  database url username
	DatabaseUser string `name:"DatabaseUser"`
	// will replace  database url password
	DatabasePassword string `name:"DatabasePassword"`
	// CacheURL is the url of the cache(redis).
	CacheURL string `name:"CacheURL"`
	// will replace  cache url username
	CacheUser string `name:"CacheUser"`
	// will replace  cache url password
	CachePassword string `name:"CachePassword"`
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

func (sr *SettingsResult) loadFromEnv() error {
	return utility.Load(sr)
}

// CreateSettings load orm settings from environment
func CreateSettings() (SettingsResult, error) {
	var out SettingsResult
	err := out.loadFromEnv()
	return out, err
}

// SettingsModule is a module that provides the settings.
var SettingsModule = fx.Provide(
	func() (SettingsResult, error) {
		return CreateSettings()
	},
)
