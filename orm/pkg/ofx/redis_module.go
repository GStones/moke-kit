package ofx

import (
	"context"

	goredis "github.com/go-redis/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nerrors"
)

// RedisParams provides the RedisParams to the mfx dependency graph.
type RedisParams struct {
	fx.In
	Redis *goredis.Client `name:"Redis"`
	Cache *goredis.Client `name:"Cache"`
}

// RedisResult provides the RedisResult to the mfx dependency graph.
type RedisResult struct {
	fx.Out
	Redis *goredis.Client `name:"Redis"`
	Cache *goredis.Client `name:"Cache"`
}

func (rr *RedisResult) Execute(
	lc fx.Lifecycle,
	l *zap.Logger,
	n SettingsParams,
) (err error) {
	if n.RedisURL != "" {
		l.Info("Connecting to redis", zap.String("host", n.RedisURL))
		rr.Redis = goredis.NewClient(&goredis.Options{
			Addr:     n.RedisURL,
			Password: n.RedisPassword,
			DB:       0,
		})
		rr.Cache = goredis.NewClient(&goredis.Options{
			Addr:     n.RedisURL,
			Password: n.RedisPassword,
			DB:       1,
		})

	} else {
		return nerrors.ErrMissingNosqlURL
	}

	if rr.Redis != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				_ = rr.Redis.Close()
				return rr.Cache.Close()
			},
		})
	}
	return
}

// RedisModule is the module for redis driver
// github.com/go-redis/redis
var RedisModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		n SettingsParams,
	) (out RedisResult, err error) {
		err = out.Execute(lc, l, n)
		return
	},
)
