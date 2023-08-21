package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility/uconfig"
)

type SecuritySettingsParams struct {
	fx.In

	TlsCert string `name:"TlsCert"`
	TlsKey  string `name:"TlsKey"`
	Secure  bool   `name:"Secure"`
}

type SecuritySettingsResult struct {
	fx.Out

	TlsCert string `name:"TlsCert" envconfig:"TLS_CERT"`
	TlsKey  string `name:"TlsKey" envconfig:"TLS_KEY"`
	Secure  bool   `name:"Secure" envonfig:"SECURE"`
}

func (g *SecuritySettingsResult) LoadFromEnv() (err error) {
	err = uconfig.Load(g)
	return
}

var SecuritySettingsModule = fx.Provide(
	func() (out SecuritySettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
