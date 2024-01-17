package agonesfx

import (
	allocation "agones.dev/agones/pkg/allocation/go"
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/3rd/agones/internal/allocate"
	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/utility"
)

// Agones allocate client module
// https://agones.dev/site/docs/advanced/allocator-service/

type AllocateParams struct {
	fx.In

	AllocateClient allocation.AllocationServiceClient `name:"AllocateClient"`
}

type AllocateResult struct {
	fx.Out

	AllocateClient allocation.AllocationServiceClient `name:"AllocateClient"`
}

// NewAllocateClient creates a new AllocateClient, requires a host and security settings.
// Agones need security settings(mTls) to connect to the allocator server.
// https://agones.dev/site/docs/advanced/allocator-service/#client-certificate
func NewAllocateClient(sSetting AllocateSettingsParams) (allocation.AllocationServiceClient, error) {
	if conn, err := utility.DialWithSecurity(
		sSetting.AllocateServiceUrl,
		sSetting.ClientCert,
		sSetting.ClientKey,
		"",
		sSetting.ServerCaCert,
	); err != nil {
		return nil, err
	} else {
		return allocation.NewAllocationServiceClient(conn), nil
	}
}

func NewAllocateClientMock(url string) (allocation.AllocationServiceClient, error) {
	return &allocate.MockAllocationServiceClient{URL: url}, nil
}

// AllocateClientModule is a fx module that provides an AllocateClient
var AllocateClientModule = fx.Provide(
	func(
		params mfx.AppParams,
		sSetting AllocateSettingsParams,
	) (out AllocateResult, err error) {
		if utility.ParseDeployments(params.Deployment).IsProd() {
			if cli, e := NewAllocateClient(sSetting); e != nil {
				err = e
			} else {
				out.AllocateClient = cli
			}
		} else {
			if cli, e := NewAllocateClientMock(sSetting.MockAllocateUrl); err != nil {
				err = e
			} else {
				out.AllocateClient = cli
			}
		}
		return
	},
)
