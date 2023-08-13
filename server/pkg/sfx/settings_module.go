package sfx

import (
	"go.uber.org/fx"
	"moke-kit/utility/uconfig"
)

type SettingsParams struct {
	fx.In

	Port        int32  `name:"Port"`
	GatewayHost string `name:"GatewayHost"`
	Mod         string `name:"Mod"`
	ZinxTcpPort int32  `name:"ZinxTcpPort"`
	ZinxWSPort  int32  `name:"ZinxWSPort"`
}

// SettingsResult loads from the environment and its members are injected into the tfx dependency graph.
type SettingsResult struct {
	fx.Out

	Port        int32  `name:"Port"  envconfig:"PORT" default:"8081"`
	GatewayHost string `name:"GatewayHost" envconfig:"GATEWAY_HOST" default:"localhost:8081"`
	Mod         string `name:"Mod" envconfig:"ZINX_MOD" default:"all" help:"grpc,tcp,websocket,http,all" `
	ZinxTcpPort int32  `name:"ZinxTcpPort" envconfig:"ZINX_TCP_PORT" default:"8888"`
	ZinxWSPort  int32  `name:"ZinxWSPort" envconfig:"ZINX_WS_PORT" default:"8889"`
}

func (g *SettingsResult) LoadFromEnv() (err error) {
	err = uconfig.Load(g)
	return
}

var SettingsModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)