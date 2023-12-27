package ofx

import (
	"context"
	"net/url"

	goredis "github.com/redis/go-redis/v9"

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
) error {
	if n.CacheURL == "" {
		return nil
	} else if u, e := url.Parse(n.CacheURL); e != nil {
		return e
	} else if u.Scheme != "redis" {
		l.Error("Invalid redis url", zap.String("url", n.CacheURL))
		return nerrors.ErrInvalidNosqlURL
	} else {
		username := u.User.Username()
		password, _ := u.User.Password()
		rr.Redis = goredis.NewClient(&goredis.Options{
			Addr:     u.Host,
			Username: username,
			Password: password,
			DB:       0,
		})
		rr.Cache = goredis.NewClient(&goredis.Options{
			Addr:     u.Host,
			Username: username,
			Password: password,
			DB:       1,
		})
	}
	l.Info("Connecting to redis", zap.String("host", n.CacheURL))

	if rr.Redis != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				_ = rr.Redis.Close()
				return rr.Cache.Close()
			},
		})
	}
	return nil
}

// RedisModule is the module for redis driver
// github.com/redis/go-redis/v9
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
