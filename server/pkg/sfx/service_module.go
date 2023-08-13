package sfx

import (
	"go.uber.org/fx"

	"moke-kit/server/siface"
)

type GrpcServiceParams struct {
	fx.In

	GrpcService siface.IGrpcService `group:"GrpcService"`
}

type GrpcServiceResult struct {
	fx.Out

	GrpcService siface.IGrpcService `group:"GrpcService"`
}

type ZinxServiceParams struct {
	fx.In

	ZinxService siface.IZinxService `group:"ZinxService"`
}

type ZinxServiceResult struct {
	fx.Out

	ZinxService siface.IZinxService `group:"ZinxService"`
}

type GatewayServiceParams struct {
	fx.In

	GatewayService siface.IGatewayService `group:"GatewayService"`
}

type GatewayServiceResult struct {
	fx.Out

	GatewayService siface.IGatewayService `group:"GatewayService"`
}
