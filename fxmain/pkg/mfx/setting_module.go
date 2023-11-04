package mfx

import (
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

	AppName    string `name:"AppName" envconfig:"APP_NAME" default:"app" ignored:"true" `
	AppId      string `name:"AppId" envconfig:"APP_ID" default:"app"`
	Deployment string `name:"Deployment" envconfig:"DEPLOYMENT" default:"local"`
	Version    string `name:"Version" default:"0.0.2"`
}

func (ar *AppResult) LoadConstant(value string) error {
	ar.AppName = value
	return nil
}

func (ar *AppResult) LoadFromEnv() (err error) {
	err = utility.Load(ar)
	return
}

var SettingModule = fx.Provide(
	func() (out AppResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
