package main

import (
	"github.com/gstones/moke-kit/demo/internal/demo"
	"github.com/gstones/moke-kit/demo/pkg/dfx"
	"github.com/gstones/moke-kit/fxmain"
	"github.com/gstones/moke-kit/mq/pkg/qfx"
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
