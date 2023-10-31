package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/key"
)

var (
	// ExpireRangeMin is the minimum expire time
	ExpireRangeMin = 40 * time.Minute
	// ExpireRangeMax is the maximum expire time
	ExpireRangeMax = 60 * time.Minute
)

// RedisCache is a redis cache
type RedisCache struct {
	logger *zap.Logger
	*redis.Client
}

// CreateRedisCache creates a redis cache
func CreateRedisCache(logger *zap.Logger, client *redis.Client) *RedisCache {
	return &RedisCache{logger, client}
}

// GetCache gets cache
func (c *RedisCache) GetCache(key key.Key, doc any) bool {
	if res := c.Get(context.Background(), key.String()); res.Err() != nil {
		return false
	} else if data, err := res.Bytes(); err != nil {
		return false
	} else if err := json.Unmarshal(data, doc); err != nil {
		return false
	}
	return true
}

// SetCache sets cache
func (c *RedisCache) SetCache(key key.Key, doc any) {
	expire := random.RandInt(int(ExpireRangeMin), int(ExpireRangeMax))
	if data, err := json.Marshal(doc); err != nil {
		return
	} else if res := c.Set(context.Background(), key.String(), data, time.Duration(expire)); res.Err() != nil {
		return
	}
}

// DeleteCache deletes cache
func (c *RedisCache) DeleteCache(key key.Key) {
	if res := c.Del(context.Background(), key.String()); res.Err() != nil {
		return
	}
}
