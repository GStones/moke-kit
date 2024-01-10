package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/server/siface"
)

type GrpcServiceParams struct {
	fx.In

	GrpcServices []siface.IGrpcService `group:"GrpcService"`
}

type GrpcServiceResult struct {
	fx.Out

	GrpcService siface.IGrpcService `group:"GrpcService"`
}

type ZinxServiceParams struct {
	fx.In

	ZinxServices []siface.IZinxService `group:"ZinxService"`
}

type ZinxServiceResult struct {
	fx.Out

	ZinxService siface.IZinxService `group:"ZinxService"`
}

type GatewayServiceParams struct {
	fx.In

	GatewayServices []siface.IGatewayService `group:"GatewayService"`
}

type GatewayServiceResult struct {
	fx.Out

	GatewayService siface.IGatewayService `group:"GatewayService"`
}
