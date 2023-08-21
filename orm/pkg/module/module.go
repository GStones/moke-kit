package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/orm/pkg/nfx"
)

var Module = fx.Module("nosql", fx.Options(
	nfx.SettingsModule,
	nfx.MongoPureModule,
	nfx.DocumentStoreModule,
	nfx.RedisModule,
	nfx.GormModule,
))
