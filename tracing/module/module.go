package module

import (
	"go.uber.org/fx"

	"moke-kit/tracing/tfx"
)

var Module = fx.Module("tracing", fx.Options(
	tfx.SettingsModule,
	tfx.TracerModule,
))
