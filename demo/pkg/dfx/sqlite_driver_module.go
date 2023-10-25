package dfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"

	"github.com/gstones/moke-kit/orm/pkg/ofx"
)

var SqliteDriverModule = fx.Provide(
	func(
		l *zap.Logger,
		s SettingsParams,
	) (out ofx.GormDriverResult, err error) {
		out.Dialector = sqlite.Open(s.GormDns)
		return out, nil
	},
)
