package fx

import (
	"time"

	"go.uber.org/fx"
)

type NoSQLSettings struct {
	fx.In

	NoSqlUser        string        `name:"NoSqlUser"`
	NoSqlPassword    string        `name:"NoSqlPassword"`
	DocumentStoreUrl string        `name:"DocumentStoreUrl"`
	MemoryStoreUrl   string        `name:"MemoryStoreUrl"`
	SessionStoreName string        `name:"SessionStoreName"`
	GCInterval       time.Duration `name:"GCInterval"`
}

type NoSQLSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	NoSqlUser        string        `name:"NoSqlUser" envconfig:"NOSQL_USERNAME"`
	NoSqlPassword    string        `name:"NoSqlPassword" envconfig:"NOSQL_PASSWORD"`
	DocumentStoreUrl string        `name:"DocumentStoreUrl" envconfig:"DOCUMENT_STORE_URL"`
	MemoryStoreUrl   string        `name:"MemoryStoreUrl" envconfig:"MEMORY_STORE_URL"`
	SessionStoreName string        `name:"SessionStoreName" default:"sessions" envconfig:"SESSION_STORE"`
	GCInterval       time.Duration `name:"GCInterval" default:"5m" envconfig:"GCINTERVAL"`
}

func (g *NoSQLSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var NoSQLSettingsModule = fx.Provide(
	func() (out NoSQLSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
