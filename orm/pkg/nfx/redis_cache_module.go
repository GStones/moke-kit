package nfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/cache"
	"github.com/gstones/moke-kit/orm/nosql/diface"
)

// RedisCacheParams provides the RedisCacheParams to the mfx dependency graph.
type RedisCacheParams struct {
	fx.In
	RedisCache diface.ICache `name:"RedisCache"`
}

// RedisCacheResult provides the RedisCacheResult to the mfx dependency graph.
type RedisCacheResult struct {
	fx.Out
	RedisCache diface.ICache `name:"RedisCache"`
}

// Execute initializes the RedisCacheResult.
func (c *RedisCacheResult) Execute(
	l *zap.Logger,
	rParams RedisParams,
) (err error) {
	c.RedisCache = cache.CreateRedisCache(l, rParams.Cache)
	return
}

// RedisCacheModule provides the RedisCacheModule to the mfx dependency graph.
var RedisCacheModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		rParams RedisParams,
	) (out RedisCacheResult, err error) {
		err = out.Execute(l, rParams)
		return
	},
)
