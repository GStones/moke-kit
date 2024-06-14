package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/server/siface"
)

// All service fx struct

// GrpcServiceParams module params for injecting GrpcService
type GrpcServiceParams struct {
	fx.In

	GrpcServices []siface.IGrpcService `group:"GrpcService"`
}

// GrpcServiceResult module result for exporting GrpcService
type GrpcServiceResult struct {
	fx.Out

	GrpcService siface.IGrpcService `group:"GrpcService"`
}

// ZinxServiceParams module params for injecting ZinxService
type ZinxServiceParams struct {
	fx.In

	ZinxServices []siface.IZinxService `group:"ZinxService"`
}

// ZinxServiceResult module result for exporting ZinxService
type ZinxServiceResult struct {
	fx.Out

	ZinxService siface.IZinxService `group:"ZinxService"`
}

// GatewayServiceParams module params for injecting GatewayService
type GatewayServiceParams struct {
	fx.In

	GatewayServices []siface.IGatewayService `group:"GatewayService"`
}

// GatewayServiceResult module result for exporting GatewayService
type GatewayServiceResult struct {
	fx.Out

	GatewayService siface.IGatewayService `group:"GatewayService"`
}
