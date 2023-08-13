package dfx

import (
	"go.uber.org/fx"

	"moke-kit/utility/uconfig"
)

type SettingsParams struct {
	fx.In

	DemoUrl string `name:"DemoUrl"`
	DbName  string `name:"DbName"`
	GormDns string `name:"GormDns"`
}

type SettingsResult struct {
	fx.Out

	DemoUrl string `name:"DemoUrl" envconfig:"DEMO_URL" default:"localhost:8081"`
	DbName  string `name:"DbName" envconfig:"DB_NAME" default:"demo"`
	GormDns string `name:"GormDns" envconfig:"GORM_DNS" default:"demo.db"`
}

func (g *SettingsResult) LoadFromEnv() (err error) {
	err = uconfig.Load(g)
	return
}

var SettingsModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
