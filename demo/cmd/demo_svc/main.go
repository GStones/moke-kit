package main

import (
	"github.com/gstones/moke-kit/demo/internal/demo"
	"github.com/gstones/moke-kit/demo/pkg/dfx"
	"github.com/gstones/moke-kit/fxmain"
	"github.com/gstones/moke-kit/mq/pkg/mfx"
)

func main() {
	fxmain.Main(
		dfx.SettingsModule,
		// config sqlite db module
		//(you can customize it as gorm support db adapter)
		dfx.SqliteDriverModule,
		// nats message queue
		mfx.NatsModule,
		// auth function(optional)
		demo.AuthModule,
		// grpc server
		demo.GrpcModule,
		// http server
		demo.GatewayModule,
		// tcp/websocket server
		demo.ZinxModule,
	)
}
