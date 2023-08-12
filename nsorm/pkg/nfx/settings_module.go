package nfx

import (
	"moke-kit/utility/uconfig"
	"time"

	"go.uber.org/fx"
)

type SettingsParams struct {
	fx.In

	NoSqlUser        string        `name:"NoSqlUser"`
	NoSqlPassword    string        `name:"NoSqlPassword"`
	DocumentStoreUrl string        `name:"DocumentStoreUrl"`
	MemoryStoreUrl   string        `name:"MemoryStoreUrl"`
	SessionStoreName string        `name:"SessionStoreName"`
	GCInterval       time.Duration `name:"GCInterval"`
}

type SettingsResult struct {
	fx.Out

	NoSqlUser        string        `name:"NoSqlUser" envconfig:"NOSQL_USERNAME"`
	NoSqlPassword    string        `name:"NoSqlPassword" envconfig:"NOSQL_PASSWORD"`
	DocumentStoreUrl string        `name:"DocumentStoreUrl" envconfig:"DOCUMENT_STORE_URL" default:"mongodb://localhost:27017"`
	MemoryStoreUrl   string        `name:"MemoryStoreUrl" envconfig:"MEMORY_STORE_URL"`
	SessionStoreName string        `name:"SessionStoreName" default:"sessions" envconfig:"SESSION_STORE"`
	GCInterval       time.Duration `name:"GCInterval" default:"5m" envconfig:"GCINTERVAL"`
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
