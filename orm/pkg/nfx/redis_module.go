package nfx

import (
	"context"
	"net/url"

	goredis "github.com/go-redis/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/orm/nerrors"
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
	if n.MemoryStoreUrl != "" {
		if u, e := url.Parse(n.MemoryStoreUrl); e != nil {
			err = e
		} else {
			switch u.Scheme {
			case "redis":
				password, set := u.User.Password()
				if !set {
					password = n.NoSqlPassword
				}
				l.Info("Connecting to redis", zap.String("host", u.Host))
				rr.Redis = goredis.NewClient(&goredis.Options{
					Addr:     u.Host,
					Password: password,
					DB:       0,
				})
				rr.Cache = goredis.NewClient(&goredis.Options{
					Addr:     u.Host,
					Password: password,
					DB:       1,
				})
			default:
				return nerrors.ErrInvalidNosqlURL
			}
		}
	} else {
		return nerrors.ErrMissingNosqlURL
	}

	if rr.Redis != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				rr.Redis.Close()
				return rr.Cache.Close()
			},
		})
	}

	return
}

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
