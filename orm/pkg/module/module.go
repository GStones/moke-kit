package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/orm/pkg/nfx"
)

// Module is the orm module
// inject redis,mongo,gorm
var Module = fx.Module("orm",
	nfx.SettingsModule,
	nfx.MongoPureModule,
	nfx.DocumentStoreModule,
	nfx.RedisModule,
	nfx.GormModule,
)
