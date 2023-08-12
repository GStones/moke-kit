package main

import (
	"moke-kit/demo/internal/demo"
	"moke-kit/demo/pkg/dfx"
	"moke-kit/fxmain"
	mq "moke-kit/mq/pkg/module"
	"moke-kit/mq/pkg/qfx"
)

func main() {
	fxmain.Main(
		dfx.SettingsModule,
		// db
		dfx.DemoDBModule,
		// message queue
		mq.Module,
		qfx.NatsModule,
		// grpc server
		demo.GrpcModule,
		// http server
		demo.GatewayModule,
		// tcp/websocket server
		demo.ZinxModule,
	)
}
