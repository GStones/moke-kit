package dfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/demo/internal/demo"
	"moke-kit/server/pkg/sfx"
)

var DemoModule = fx.Provide(
	func(
		l *zap.Logger,
	) (out sfx.GrpcServiceResult, err error) {
		if s, err := demo.NewService(l); err != nil {
			return out, err
		} else {
			out.GrpcService = s
		}
		return
	},
)

var DemoGatewayModule = fx.Provide(
	func(
		l *zap.Logger,
	) (out sfx.GatewayServiceResult, err error) {
		if s, err := demo.NewService(l); err != nil {
			return out, err
		} else {
			out.GatewayService = s
		}
		return
	},
)
