package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/tracing/tfx"
)

var Module = fx.Module("tracing", fx.Options(
	tfx.SettingsModule,
	tfx.TracerModule,
))
