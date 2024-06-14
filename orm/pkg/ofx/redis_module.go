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

	//use redis db 0
	Redis *goredis.Client `name:"Redis"`
	//use redis db 1
	Cache *goredis.Client `name:"Cache"`
}

// RedisResult provides the RedisResult to the mfx dependency graph.
type RedisResult struct {
	fx.Out
	//use redis db 0
	Redis *goredis.Client `name:"Redis"`
	//use redis db 1
	Cache *goredis.Client `name:"Cache"`
}

func (rr *RedisResult) init(
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
		if n.CacheUser != "" {
			username = n.CacheUser
		}
		password, _ := u.User.Password()
		if n.CachePassword != "" {
			password = n.CachePassword
		}
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
	l.Info("Connecting redis", zap.String("host", n.CacheURL))
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

// CreateRedis creates a new Redis driver
func CreateRedis(
	lc fx.Lifecycle,
	l *zap.Logger,
	n SettingsParams,
) (RedisResult, error) {
	var out RedisResult
	err := out.init(lc, l, n)
	return out, err
}

// RedisModule is the module for redis driver
// github.com/redis/go-redis/v9
var RedisModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		n SettingsParams,
	) (RedisResult, error) {
		return CreateRedis(lc, l, n)
	},
)
