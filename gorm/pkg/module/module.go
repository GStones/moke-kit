package module

import (
	"go.uber.org/fx"

	"moke-kit/gorm/pkg/nfx"
)

var Module = fx.Module("nosql", fx.Options(
	nfx.MongoPureModule,
	nfx.DocumentStoreModule,
	nfx.SettingsModule,
	nfx.RedisModule,
))
