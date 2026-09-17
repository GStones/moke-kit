package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/orm/pkg/ofx"
)

// Module is the default orm graph: Mongo document store + Redis.
// GORM is opt-in via ofx.GormModule (no in-kit dialector).
var Module = fx.Module("orm",
	ofx.SettingsModule,
	ofx.MongoPureModule,
	ofx.DocumentStoreModule,
	ofx.RedisModule,
)
