package mfx

import (
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/utility"

	"go.uber.org/fx"
)

type AppParams struct {
	fx.In

	AppName    string `name:"AppName"`
	AppId      string `name:"AppId"`
	Deployment string `name:"Deployment"`
	Version    string `name:"Version"`
}

type AppResult struct {
	fx.Out

	AppName    string `name:"AppName" envconfig:"APP_NAME" default:"app"`
	AppId      string `name:"AppId" envconfig:"APP_ID" default:"app"`
	Deployment string `name:"Deployment" envconfig:"DEPLOYMENT" default:"local"`
	Version    string `name:"Version" envconfig:"VERSION" default:"0.0.1"`
}

func (ar *AppResult) loadFromEnv() error {
	if err := utility.Load(ar); err != nil {
		return err
	}
	key.SetNamespace(ar.Deployment)
	return nil
}

// SettingModule is a module that provides the application settings.
var SettingModule = fx.Provide(
	func() (out AppResult, err error) {
		err = out.loadFromEnv()
		return
	},
)
