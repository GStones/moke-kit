package cache

import (
	"encoding/json"
	"time"

	"github.com/duke-git/lancet/v2/random"
	"github.com/go-redis/redis"
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
	if res := c.Get(key.String()); res.Err() != nil {
		return false
	} else {
		if err := json.Unmarshal([]byte(res.Val()), &doc); err != nil {
			return false
		}
		return true
	}
}

// SetCache sets cache
func (c *RedisCache) SetCache(key key.Key, doc any) {
	if data, err := json.Marshal(doc); err != nil {
		return
	} else {
		expire := random.RandInt(int(ExpireRangeMin), int(ExpireRangeMax))
		if res := c.Set(key.String(), data, time.Duration(expire)); res.Err() != nil {
			return
		}
	}
}

// DeleteCache deletes cache
func (c *RedisCache) DeleteCache(key key.Key) {
	if res := c.Del(key.String()); res.Err() != nil {
		return
	}
}
