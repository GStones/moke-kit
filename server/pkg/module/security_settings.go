package fxsvcapp

import (
	"go.uber.org/fx"
	"moke-kit/common/config"
)

type SecuritySettingsParams struct {
	fx.In

	TlsCert string `name:"TlsCert"`
	TlsKey  string `name:"TlsKey"`

	SecureClients bool `name:"SecureClients"`
}

type SecuritySettingsResult struct {
	fx.Out

	TlsCert string `name:"TlsCert" envconfig:"TLS_CERT"`
	TlsKey  string `name:"TlsKey" envconfig:"TLS_KEY"`

	SecureClients bool `name:"SecureClients" envonfig:"SECURE_CLIENTS"`
}

func (g *SecuritySettingsResult) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var SecuritySettingsModule = fx.Provide(
	func() (out SecuritySettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
