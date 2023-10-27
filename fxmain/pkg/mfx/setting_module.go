package mfx

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	err = utility.Load(ar)
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
