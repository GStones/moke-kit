package main

import (
	"moke-kit/demo/pkg/dfx"
	fxapp "moke-kit/fxmain"
)

func main() {
	fxapp.Main(
		dfx.DemoModule,
		dfx.DemoGatewayModule,
	)
}
