package ofx

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
func (c *RedisCacheResult) init(
	l *zap.Logger,
	rParams RedisParams,
) error {
	c.RedisCache = cache.CreateRedisCache(l, rParams.Cache)
	return nil
}

// CreateRedisCache creates a redis cathe .
func CreateRedisCache(
	l *zap.Logger,
	rParams RedisParams,
) (RedisCacheResult, error) {
	var out RedisCacheResult
	err := out.init(l, rParams)
	return out, err
}

// RedisCacheModule provides the RedisCacheModule to the mfx dependency graph.
var RedisCacheModule = fx.Provide(
	func(
		l *zap.Logger,
		rParams RedisParams,
	) (RedisCacheResult, error) {
		return CreateRedisCache(l, rParams)
	},
)
