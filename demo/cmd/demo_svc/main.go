package main

import (
	"moke-kit/demo/internal/demo"
	"moke-kit/demo/pkg/dfx"
	"moke-kit/fxmain"
	"moke-kit/mq/pkg/qfx"
)

func main() {
	fxmain.Main(
		dfx.SettingsModule,
		// db
		dfx.DemoDBModule,
		// sqlite db
		dfx.SqliteDriverModule,
		// nats message queue
		qfx.NatsModule,
		// grpc server
		demo.GrpcModule,
		// http server
		demo.GatewayModule,
		// tcp/websocket server
		demo.ZinxModule,
	)
}
