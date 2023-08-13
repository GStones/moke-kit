package nfx

import (
	"go.uber.org/fx"
	"moke-kit/utility/uconfig"
)

type SettingsParams struct {
	fx.In

	NoSqlUser      string `name:"NoSqlUser"`
	NoSqlPassword  string `name:"NoSqlPassword"`
	NosqlUrl       string `name:"NosqlUrl"`
	MemoryStoreUrl string `name:"MemoryStoreUrl"`
}

type SettingsResult struct {
	fx.Out

	NoSqlUser      string `name:"NoSqlUser" envconfig:"NOSQL_USERNAME"`
	NoSqlPassword  string `name:"NoSqlPassword" envconfig:"NOSQL_PASSWORD"`
	NoSqlUrl       string `name:"NosqlUrl" envconfig:"NOSQL_URL" default:"mongodb://localhost:27017"`
	MemoryStoreUrl string `name:"MemoryStoreUrl" envconfig:"MEMORY_STORE_URL"`
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
