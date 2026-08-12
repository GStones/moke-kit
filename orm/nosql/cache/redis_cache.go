package cache

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/key"
)

var (
	// ExpireRangeMin is the minimum expire time used by callers for jittered TTLs.
	ExpireRangeMin = 40 * time.Minute
	// ExpireRangeMax is the maximum expire time used by callers for jittered TTLs.
	ExpireRangeMax = 60 * time.Minute
)

// RedisCache is a redis HASH cache.
type RedisCache struct {
	logger *zap.Logger
	*redis.Client
}

var _ diface.ICache = (*RedisCache)(nil)

// CreateRedisCache creates a redis cache.
func CreateRedisCache(logger *zap.Logger, client *redis.Client) *RedisCache {
	return &RedisCache{logger, client}
}

func isWrongType(err error) bool {
	return err != nil && strings.Contains(err.Error(), "WRONGTYPE")
}

// clearLegacyKey deletes a pre-HASH string value that conflicts with HSET/HGETALL.
func (c *RedisCache) clearLegacyKey(ctx context.Context, keyStr string, err error) {
	if !isWrongType(err) {
		return
	}
	c.logger.Warn("clearing legacy non-HASH cache key", zap.String("key", keyStr), zap.Error(err))
	_ = c.Del(ctx, keyStr).Err()
}

// GetCache loads HASH fields. When fields is empty, all fields are returned.
func (c *RedisCache) GetCache(ctx context.Context, key key.Key, fields ...string) map[string]any {
	if c == nil || c.Client == nil {
		return nil
	}

	keyStr := key.String()
	if len(fields) > 0 {
		result, err := c.HMGet(ctx, keyStr, fields...).Result()
		if err != nil {
			c.clearLegacyKey(ctx, keyStr, err)
			if !isWrongType(err) {
				c.logger.Error("get cache failed", zap.Error(err), zap.String("key", keyStr))
			}
			return nil
		}
		data := make(map[string]any, len(fields))
		for i, field := range fields {
			if result[i] != nil {
				data[field] = result[i]
			}
		}
		if len(data) == 0 {
			return nil
		}
		return data
	}

	allData, err := c.HGetAll(ctx, keyStr).Result()
	if err != nil {
		c.clearLegacyKey(ctx, keyStr, err)
		if !isWrongType(err) {
			c.logger.Error("get cache failed", zap.Error(err), zap.String("key", keyStr))
		}
		return nil
	}
	if len(allData) == 0 {
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

// SetCache writes HASH fields with HSET and optional EXPIRE.
func (c *RedisCache) SetCache(
	ctx context.Context,
	key key.Key,
	data map[string]any,
	expire time.Duration,
) error {
	if c == nil || c.Client == nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	keyStr := key.String()
	if res := c.HSet(ctx, keyStr, data); res.Err() != nil {
		if !isWrongType(res.Err()) {
			return res.Err()
		}
		// Upgrade path: older versions stored a JSON string under the same key.
		c.clearLegacyKey(ctx, keyStr, res.Err())
		if retry := c.HSet(ctx, keyStr, data); retry.Err() != nil {
			return retry.Err()
		}
	}
	if expire > 0 {
		if res := c.Expire(ctx, keyStr, expire); res.Err() != nil {
			return res.Err()
		}
	}
	return nil
}

// DeleteCache deletes cache.
func (c *RedisCache) DeleteCache(ctx context.Context, key key.Key) {
	if c == nil || c.Client == nil {
		return
	}
	if res := c.Del(ctx, key.String()); res.Err() != nil {
		return
	}
}
