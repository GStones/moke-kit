package mfx

import (
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

func (l *AppResult) LoadConstant(value string) error {
	l.AppName = value
	return nil
}

func (l *AppResult) LoadFromExecutable() (err error) {
	if exeName, e := os.Executable(); e != nil {
		err = e
	} else {
		exeName = filepath.Base(exeName)
		if runtime.GOOS == "windows" && strings.HasSuffix(exeName, ".exe") {
			exeName = exeName[:len(exeName)-4]
		}
		l.AppName = exeName
	}
	return
}

var AppModule = fx.Provide(
	func() (out AppResult, err error) {
		err = out.LoadFromExecutable()
		return
	},
)
