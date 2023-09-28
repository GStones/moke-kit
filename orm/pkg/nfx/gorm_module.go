package nfx

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormParams struct {
	fx.In

	GormDB *gorm.DB `name:"GormDB"`
}

type GormResult struct {
	fx.Out

	GormDB *gorm.DB `name:"GormDB"`
}

type GormDriverParams struct {
	fx.In

	Dialector gorm.Dialector `name:"Dialector"`
}

type GormDriverResult struct {
	fx.Out

	Dialector gorm.Dialector `name:"Dialector"`
}

func (mr *GormResult) NewDocument(
	lc fx.Lifecycle,
	logger *zap.Logger,
	dialector gorm.Dialector,
) (err error) {
	if dialector == nil {
		logger.Info("no gorm driver")
		return
	} else if db, err := gorm.Open(dialector, &gorm.Config{}); err != nil {
		return err
	} else {
		logger.Info("open gorm driver", zap.String("dialector", dialector.Name()))
		mr.GormDB = db

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("close gorm driver", zap.String("dialector", dialector.Name()))
				sqlDB, err := db.DB()
				if err != nil {
					return err
				}
				return sqlDB.Close()
			},
		})

	}
	return
}

// GormModule is the module for gorm driver
// https://github.com/go-gorm/gorm
var GormModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		dParams GormDriverParams,
	) (dOut GormResult, err error) {
		err = dOut.NewDocument(lc, l, dParams.Dialector)
		return
	},
)
