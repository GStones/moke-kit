package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

type SettingsParams struct {
	fx.In

	Port        int32 `name:"Port"`
	ZinxTcpPort int32 `name:"ZinxTcpPort"`
	ZinxWSPort  int32 `name:"ZinxWSPort"`
	Timeout     int32 `name:"Timeout"`
	RateLimit   int32 `name:"RateLimit"`
}

// SettingsResult loads from the environment and its members are injected into the tfx dependency graph.
type SettingsResult struct {
	fx.Out

	Port        int32 `name:"Port"  envconfig:"PORT" default:"8081"`
	ZinxTcpPort int32 `name:"ZinxTcpPort" envconfig:"ZINX_TCP_PORT" default:"8888"`
	ZinxWSPort  int32 `name:"ZinxWSPort" envconfig:"ZINX_WS_PORT" default:""`
	Timeout     int32 `name:"Timeout" envconfig:"TIMEOUT" default:"10"`        // connection heartbeat timeout
	RateLimit   int32 `name:"RateLimit" envconfig:"RATE_LIMIT" default:"1000"` // rate limit per second
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
