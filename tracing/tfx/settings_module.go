package tfx

import (
	"go.uber.org/fx"
	"moke-kit/utility/uconfig"
)

type SettingsParams struct {
	fx.In

	TraceProvider    string   `name:"TraceProvider"`
	TraceAgentHost   string   `name:"TraceAgentHost"`
	TraceAgentPort   int      `name:"TraceAgentPort"`
	TraceServiceName string   `name:"TraceServiceName"`
	TraceTags        []string `name:"TraceTags"`
}

type SettingsResult struct {
	fx.Out

	TraceProvider    string   `name:"TraceProvider" envconfig:"TRACE_PROVIDER"`
	TraceAgentHost   string   `name:"TraceAgentHost" envconfig:"TRACE_AGENT_HOST"`
	TraceAgentPort   int      `name:"TraceAgentPort" envconfig:"TRACE_AGENT_PORT"`
	TraceServiceName string   `name:"TraceServiceName" envconfig:"TRACE_SERVICE_NAME"`
	TraceTags        []string `name:"TraceTags" envconfig:"TRACE_TAGS"`
}

func (l *SettingsResult) LoadFromEnv() (err error) {
	err = uconfig.Load(l)
	return
}

var SettingsModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
