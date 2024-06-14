package ofx

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

func (mr *GormResult) init(
	lc fx.Lifecycle,
	logger *zap.Logger,
	dialector gorm.Dialector,
) error {
	if dialector == nil {
		logger.Info("no gorm driver")
		return nil
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
	return nil
}

// CreateGormDriver creates a new gorm driver
func CreateGormDriver(
	lc fx.Lifecycle,
	logger *zap.Logger,
	dialector gorm.Dialector,
) (GormResult, error) {
	var mr GormResult
	if err := mr.init(lc, logger, dialector); err != nil {
		return mr, err
	}
	return mr, nil
}

// GormModule is the module for gorm driver
// https://github.com/go-gorm/gorm
var GormModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		dParams GormDriverParams,
	) (GormResult, error) {
		return CreateGormDriver(lc, l, dParams.Dialector)
	},
)
