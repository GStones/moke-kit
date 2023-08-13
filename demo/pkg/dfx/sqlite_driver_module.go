package dfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"moke-kit/orm/pkg/nfx"
)

var SqliteDriverModule = fx.Provide(
	func(
		l *zap.Logger,
		s SettingsParams,
	) (out nfx.GormDriverResult, err error) {
		out.Dialector = sqlite.Open(s.GormDns)
		return out, nil
	},
)