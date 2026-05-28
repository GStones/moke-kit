package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
)

// RedisCache is a redis cache
type RedisCache struct {
	logger *zap.Logger
	*redis.Client
}

var _ diface.ICache = (*RedisCache)(nil)

// CreateRedisCache creates a redis cache
func CreateRedisCache(logger *zap.Logger, client *redis.Client) *RedisCache {
	return &RedisCache{logger, client}
}

// GetCache 获取缓存数据
func (c *RedisCache) GetCache(ctx context.Context, key key.Key, fields ...string) map[string]any {
	keyStr := key.String()
	// 根据是否指定字段使用不同的获取策略
	if len(fields) > 0 {
		// 获取指定字段
		result, err := c.HMGet(ctx, keyStr, fields...).Result()
		if err != nil {
			c.logger.Error("获取缓存失败", zap.Error(err), zap.String("key", keyStr))
			return nil
		}

		// 构建返回结果
		data := make(map[string]any, len(fields))
		for i, field := range fields {
			if result[i] != nil {
				data[field] = result[i]
			}
		}
		return data
	}

	// 获取所有字段
	allData, err := c.HGetAll(ctx, keyStr).Result()
	if err != nil {
		c.logger.Error("获取缓存失败", zap.Error(err), zap.String("key", keyStr))
		return nil
	}
	res := make(map[string]any, len(allData))
	for k, v := range allData {
		if v != "" {
			res[k] = v
		}
	}
	return res
}

// SetCache set cache with HSET
func (c *RedisCache) SetCache(
	ctx context.Context,
	key key.Key, data map[string]any,
	expire time.Duration) error {
	if res := c.HSet(ctx, key.String(), data); res.Err() != nil {
		return res.Err()
	}
	if expire > 0 {
		if res := c.Expire(ctx, key.String(), expire); res.Err() != nil {
			return res.Err()
		}
	}
	return nil
}

// DeleteCache deletes cache
func (c *RedisCache) DeleteCache(ctx context.Context, key key.Key) {
	if res := c.Del(ctx, key.String()); res.Err() != nil {
		return
	}
}
