package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

// All server settings module

type SettingsParams struct {
	fx.In

	Port        int32 `name:"Port"`        // grpc/http port
	ZinxTcpPort int32 `name:"ZinxTcpPort"` // tcp port
	ZinxWSPort  int32 `name:"ZinxWSPort"`  // websocket port
	Timeout     int32 `name:"Timeout"`     // tcp service heartbeat timeout
	RateLimit   int32 `name:"RateLimit"`   // all server type rate limit per second
}

// SettingsResult loads from the environment and its members are injected into the tfx dependency graph.
type SettingsResult struct {
	fx.Out

	Port        int32 `name:"Port"  envconfig:"PORT" default:"8081"`
	ZinxTcpPort int32 `name:"ZinxTcpPort" envconfig:"ZINX_TCP_PORT" default:"8888"`
	ZinxWSPort  int32 `name:"ZinxWSPort" envconfig:"ZINX_WS_PORT" default:""`
	Timeout     int32 `name:"Timeout" envconfig:"TIMEOUT" default:"10"`
	RateLimit   int32 `name:"RateLimit" envconfig:"RATE_LIMIT" default:"1000"`
}

func (g *SettingsResult) LoadFromEnv() (err error) {
	err = utility.Load(g)
	return
}

var SettingsModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
