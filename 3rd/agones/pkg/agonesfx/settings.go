package agonesfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

// Agones allocate client Settings
// https://agones.dev/site/docs/advanced/allocator-service/#send-allocation-request

type AgonesSettingsParams struct {
	fx.In

	// agones deployment (local/dev/prod)
	AgonesDeployment string `name:"AgonesDeployment"`
	// mock allocate url(only for non-prod deployment)
	MockAllocateUrl []string `name:"MockAllocateUrl"`
	// allocate service url(only for prod deployment)
	AllocateServiceUrl string `name:"AllocateServiceUrl"`
	ClientCert         string `name:"AllocateClientCert"`
	ClientKey          string `name:"AllocateClientKey"`
	ServerCaCert       string `name:"AllocateServerCaCert"`
}

type AgonesSettingsResult struct {
	fx.Out
	// agones deployment (local/dev/prod)
	AgonesDeployment string `name:"AgonesDeployment" envconfig:"AGONES_DEPLOYMENT" default:"local"`
	// mock allocate url(only for non-prod deployment)
	MockAllocateUrl []string `name:"MockAllocateUrl" envconfig:"MOCK_ALLOCATE_URL" default:"localhost:8888"`
	// allocate service url(only for prod deployment)
	AllocateServiceUrl string `name:"AllocateServiceUrl" envconfig:"ALLOCATE_SERVICE_URL"  default:""`
	ClientCert         string `name:"AllocateClientCert" envconfig:"ALLOCATE_CLIENT_CERT" default:"./configs/agones/tls.crt"`
	ClientKey          string `name:"AllocateClientKey" envconfig:"ALLOCATE_CLIENT_KEY" default:"./configs/agones/tls.key"`
	ServerCaCert       string `name:"AllocateServerCaCert" envconfig:"ALLOCATE_SERVER_CA_CERT" default:"./configs/agones/ca.crt"`
}

func (as *AgonesSettingsResult) loadFromEnv() error {
	return utility.Load(as)
}

// CreateAgonesSettingsModule creates an AgonesSettingsModule
func CreateAgonesSettingsModule() (AgonesSettingsResult, error) {
	out := AgonesSettingsResult{}
	err := out.loadFromEnv()
	return out, err
}

// SettingsModule is a fx module that provides an AgonesSettingsModule
var SettingsModule = fx.Provide(
	func() (out AgonesSettingsResult, err error) {
		return CreateAgonesSettingsModule()
	},
)
