package main

import (
	"moke-kit/demo/internal/demo"
	"moke-kit/demo/pkg/dfx"
	"moke-kit/fxmain"
)

func main() {
	fxmain.Main(
		dfx.SettingsModule,
		demo.DemoModule,
		demo.DemoGatewayModule,
		dfx.DemoDBModule,
	)
}
