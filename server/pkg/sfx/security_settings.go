package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

type SecuritySettingsParams struct {
	fx.In

	// client mTLS settings
	ClientCaCert string `name:"ClientCaCert"`
	ClientCert   string `name:"ClientCert"`
	ClientKey    string `name:"ClientKey"`

	// server mTLS settings
	ServerCaCert string `name:"ServerCaCert"`
	ServerCert   string `name:"ServerCert"`
	ServerKey    string `name:"ServerKey"`
	ServerName   string `name:"ServerName"`

	Secure bool `name:"Secure"`
}

type SecuritySettingsResult struct {
	fx.Out

	//client mTLS settings
	ClientCaCert string `name:"ClientCaCert" envconfig:"CLIENT_CA_CERT" default:"./configs/x509/client_ca_cert.pem"`
	ClientCert   string `name:"ClientCert" envconfig:"CLIENT_CERT" default:"./configs/x509/client_cert.pem"`
	ClientKey    string `name:"ClientKey" envconfig:"CLIENT_KEY" default:"./configs/x509/client_key.pem"`

	//server mTLS settings
	ServerCACert string `name:"ServerCaCert" envconfig:"SERVER_CA_CERT" default:"./configs/x509/ca_cert.pem"`
	ServerCert   string `name:"ServerCert" envconfig:"SERVER_CERT" default:"./configs/x509/server_cert.pem"`
	ServerKey    string `name:"ServerKey" envconfig:"SERVER_KEY" default:"./configs/x509/server_key.pem"`
	ServerName   string `name:"ServerName" envconfig:"SERVER_NAME" default:"localhost"`

	Secure bool `name:"Secure" envonfig:"SECURE" default:"false" `
}

func (g *SecuritySettingsResult) LoadFromEnv() (err error) {
	err = utility.Load(g)
	return
}

var SecuritySettingsModule = fx.Provide(
	func() (out SecuritySettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
