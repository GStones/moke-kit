package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/orm/pkg/ofx"
)

// Module is the orm module
// inject redis,mongo,gorm
var Module = fx.Module("orm",
	ofx.SettingsModule,
	ofx.MongoPureModule,
	ofx.DocumentStoreModule,
	ofx.RedisModule,
	ofx.GormModule,
)
