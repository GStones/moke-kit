package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/server/pkg/sfx"
)

var Module = fx.Module("server", fx.Options(
	sfx.ConnectionMuxModule,
	sfx.SecuritySettingsModule,
	sfx.ServersModule,
	sfx.SettingsModule,
))
