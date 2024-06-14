package mfx

import (
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/utility"

	"go.uber.org/fx"
)

// AppParams is the parameters for the application.
type AppParams struct {
	fx.In

	AppName    string `name:"AppName"`
	AppId      string `name:"AppId"`
	Deployment string `name:"Deployment"`
	Version    string `name:"Version"`
}

// AppResult is the result of the application.
type AppResult struct {
	fx.Out
	// AppName is the name of the application.
	AppName string `name:"AppName" envconfig:"APP_NAME" default:"app"`
	// AppId is the id of the application.
	AppId string `name:"AppId" envconfig:"APP_ID" default:"app"`
	// Deployment is the deployment of the application: local, dev, prod
	// you can customize it as your need local_<name> = local, dev_<name> = dev, prod_<name> = prod
	Deployment string `name:"Deployment" envconfig:"DEPLOYMENT" default:"local"`
	// Version is the version of the application.
	Version string `name:"Version" envconfig:"VERSION" default:"0.0.1"`
}

func (ar *AppResult) loadFromEnv() error {
	if err := utility.Load(ar); err != nil {
		return err
	}
	key.SetNamespace(ar.Deployment)
	return nil
}

func CreateAppModule() (AppResult, error) {
	out := AppResult{}
	err := out.loadFromEnv()
	return out, err
}

// SettingModule is a module that provides the application settings.
var SettingModule = fx.Provide(
	func() (out AppResult, err error) {
		return CreateAppModule()
	},
)
