package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/server/pkg/sfx"
)

var Module = fx.Module("server",
	sfx.SecuritySettingsModule,
	sfx.SettingsModule,
	sfx.ConnectionMuxModule,
	sfx.OTelModule,
)
