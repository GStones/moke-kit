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
		demo.Module,
		demo.GatewayModule,
		dfx.DemoDBModule,
		qfx.NatsModule,
	)
}
