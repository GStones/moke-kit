package module

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/pkg/sfx"
)

// Module is the module for the server. Binding runs here so thin apps can omit
// this module entirely (no port mux) and still use fxmain.Core.
var Module = fx.Module("server",
	sfx.SecuritySettingsModule,
	sfx.SettingsModule,
	sfx.ConnectionMuxModule,
	sfx.OTelModule,
	fx.Invoke(func(l *zap.Logger, lc fx.Lifecycle, sb ServiceBinder) error {
		return sb.Bind(l, lc)
	}),
)
