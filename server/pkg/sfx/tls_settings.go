package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

// Client&Server mTls settings module

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

	// Zero trust security model: all services must be mTLS enabled
	// if true, enable imports client for grpc/http(cmux service) clients
	MTLSEnable bool `name:"MTLSEnable"`

	// if true, enable service tls for grpc/http(cmux service) services
	TLSEnable bool `name:"TLSEnable"`
	// if true, enable Tls for tcp services(zinx service)
	TCPTlsEnable bool `name:"TCPTlsEnable"`
}

type SecuritySettingsResult struct {
	fx.Out

	//client mTLS settings
	ClientCaCert string `name:"ClientCaCert" envconfig:"CLIENT_CA_CERT" default:"./configs/tls-client/ca.crt"`
	ClientCert   string `name:"ClientCert" envconfig:"CLIENT_CERT" default:"./configs/tls-client/tls.crt"`
	ClientKey    string `name:"ClientKey" envconfig:"CLIENT_KEY" default:"./configs/tls-client/tls.key"`

	//server mTLS settings
	ServerCACert string `name:"ServerCaCert" envconfig:"SERVER_CA_CERT" default:"./configs/tls-server/ca.crt"`
	ServerCert   string `name:"ServerCert" envconfig:"SERVER_CERT" default:"./configs/tls-server/tls.crt"`
	ServerKey    string `name:"ServerKey" envconfig:"SERVER_KEY" default:"./configs/tls-server/tls.key"`
	ServerName   string `name:"ServerName" envconfig:"SERVER_NAME" default:""`

	// if true, enable mTLS for grpc/http(cmux service) services
	// Zero trust security model: all services must be mTLS enabled
	MTLSEnable bool `name:"MTLSEnable" envconfig:"MTLS_ENABLE" default:"false"`
	// if true, enable service tls for grpc/http(cmux service) services
	TLSEnable bool `name:"TLSEnable" envconfig:"TLS_ENABLE" default:"false"`
	// if true, enable Tls for tcp services(zinx service)
	TcpTlsEnable bool `name:"TCPTlsEnable" envconfig:"TCP_TLS_ENABLE" default:"false"`
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
