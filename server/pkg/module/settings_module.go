package fxsvcapp

import (
	"go.uber.org/fx"
	"moke-kit/common/config"
)

type SettingsParams struct {
	fx.In

	AppTestMode bool   `name:"AppTestMode"`
	Version     string `name:"Version"`
	Deployment  string `name:"Deployment"`
	AppId       string `name:"AppId"`
	Port        int32  `name:"Port"`
}

// SettingsResult loads from the environment and its members are injected into the fx dependency graph.
type SettingsResult struct {
	fx.Out

	AppTestMode bool   `name:"AppTestMode" ignored:"true"`
	Version     string `name:"Version" default:"unknown"`
	Deployment  string `name:"Deployment" default:"local" envconfig:"DEPLOYMENT"`
	AppId       string `name:"AppId" envconfig:"APP_ID" default:"app_id"`
	Port        int32  `name:"Port"`
}

func (g *SettingsResult) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var SettingsModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
