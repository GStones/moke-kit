package mfx

import (
	"moke-kit/utility/uconfig"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/fx"
)

type AppParams struct {
	fx.In

	AppName     string `name:"AppName"`
	AppId       string `name:"AppId"`
	Deployment  string `name:"Deployment"`
	AppTestMode bool   `name:"AppTestMode"`
	Version     string `name:"Version"`
}

type AppResult struct {
	fx.Out

	AppName     string `name:"AppName" envconfig:"APP_NAME" default:"fxapp" ignored:"true" `
	AppId       string `name:"AppId" envconfig:"APP_ID" default:"app_id"`
	Deployment  string `name:"Deployment" envconfig:"DEPLOYMENT" default:"local"`
	AppTestMode bool   `name:"AppTestMode" ignored:"true"`
	Version     string `name:"Version" default:"unknown"`
}

func (ar *AppResult) LoadConstant(value string) error {
	ar.AppName = value
	return nil
}

func (ar *AppResult) LoadFromExecutable() (err error) {
	if exeName, e := os.Executable(); e != nil {
		err = e
	} else {
		exeName = filepath.Base(exeName)
		if runtime.GOOS == "windows" && strings.HasSuffix(exeName, ".exe") {
			exeName = exeName[:len(exeName)-4]
		}
		ar.AppName = exeName
	}
	return
}

func (ar *AppResult) LoadFromEnv() (err error) {
	err = uconfig.Load(ar)
	return
}

var SettingModule = fx.Provide(
	func() (out AppResult, err error) {
		err = out.LoadFromExecutable()
		if err != nil {
			return
		}
		err = out.LoadFromEnv()
		return
	},
)
