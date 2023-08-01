package module

import (
	"go.uber.org/fx"

	"moke-kit/nosql/pkg/nfx"
)

var Module = fx.Module("nosql", fx.Options(
	nfx.DocumentStoreModule,
	nfx.SettingsModule,
	nfx.RedisModule,
))
