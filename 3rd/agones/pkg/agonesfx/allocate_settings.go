package agonesfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

// Agones allocate client Settings
// https://agones.dev/site/docs/advanced/allocator-service/#send-allocation-request

type AllocateSettingsParams struct {
	fx.In

	AllocateServiceUrl string `name:"AllocateServiceUrl"`
	ClientCert         string `name:"AllocateClientCert"`
	ClientKey          string `name:"AllocateClientKey"`
	ServerCaCert       string `name:"AllocateServerCaCert"`

	MockAllocateUrl string `name:"MockAllocateUrl"`
}

type AllocateSettingsResult struct {
	fx.Out

	AllocateServiceUrl string `name:"AllocateServiceUrl" envconfig:"ALLOCATE_SERVICE_URL"  default:""`
	ClientCert         string `name:"AllocateClientCert" envconfig:"ALLOCATE_CLIENT_CERT" default:"./configs/x509/agones/tls.crt"`
	ClientKey          string `name:"AllocateClientKey" envconfig:"ALLOCATE_CLIENT_KEY" default:"./configs/x509/agones/tls.key"`
	ServerCaCert       string `name:"AllocateServerCaCert" envconfig:"ALLOCATE_SERVER_CA_CERT" default:"./configs/x509/agones/ca.crt"`

	MockAllocateUrl string `name:"MockAllocateUrl" envconfig:"MOCK_ALLOCATE_URL" default:"localhost:8888"`
}

func (g *AllocateSettingsResult) LoadFromEnv() (err error) {
	err = utility.Load(g)
	return
}

var AllocateSettingsModule = fx.Provide(
	func() (out AllocateSettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
