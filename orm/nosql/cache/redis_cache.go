package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

var (
	// ExpireRangeMin is the minimum expire time
	ExpireRangeMin = 40 * time.Minute
	// ExpireRangeMax is the maximum expire time
	ExpireRangeMax = 60 * time.Minute
)

const compareAndSwapLua = `
local current = redis.call("GET", KEYS[2])
if not current then
	return redis.error_reply("ERR_VERSION_MISMATCH")
end
if tostring(current) ~= ARGV[1] then
	return redis.error_reply("ERR_VERSION_MISMATCH")
end
redis.call("SET", KEYS[1], ARGV[3], "PX", ARGV[4])
redis.call("SET", KEYS[2], ARGV[2], "PX", ARGV[4])
return ARGV[2]
`

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
func (c *RedisCache) GetCache(ctx context.Context, key key.Key, doc any) bool {
	if res := c.Get(ctx, key.String()); res.Err() != nil {
		return false
	} else if data, err := res.Bytes(); err != nil {
		return false
	} else if err := json.Unmarshal(data, doc); err != nil {
		return false
	}
	return true
}

// SetCache sets cache
func (c *RedisCache) SetCache(ctx context.Context, key key.Key, doc any, expire time.Duration) {
	if data, err := json.Marshal(doc); err != nil {
		return
	} else if res := c.Set(ctx, key.String(), data, expire); res.Err() != nil {
		return
	} else if version, ok := extractVersion(doc); ok {
		if res := c.Set(ctx, versionKey(key), version, expire); res.Err() != nil {
			return
		}
	}
}

// CompareAndSwapCache updates cache atomically by version with Redis Lua.
func (c *RedisCache) CompareAndSwapCache(
	ctx context.Context,
	key key.Key,
	expectedVersion noptions.Version,
	doc any,
	expire time.Duration,
) (noptions.Version, error) {
	if expire <= 0 {
		expire = ExpireRangeMin
	}
	nextVersion, ok := extractVersion(doc)
	if !ok {
		return noptions.NoVersion, fmt.Errorf("missing version in cache payload")
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return noptions.NoVersion, err
	}
	res := c.Eval(
		ctx,
		compareAndSwapLua,
		[]string{key.String(), versionKey(key)},
		expectedVersion,
		nextVersion,
		data,
		expire.Milliseconds(),
	)
	if err := res.Err(); err != nil {
		return noptions.NoVersion, err
	}
	return nextVersion, nil
}

// DeleteCache deletes cache
func (c *RedisCache) DeleteCache(ctx context.Context, key key.Key) {
	if res := c.Del(ctx, key.String(), versionKey(key)); res.Err() != nil {
		return
	}
}

func versionKey(k key.Key) string {
	return k.String() + ":version"
}

func extractVersion(doc any) (noptions.Version, bool) {
	if doc == nil {
		return noptions.NoVersion, false
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return noptions.NoVersion, false
	}
	var payload struct {
		Version noptions.Version `json:"Version"`
	}
	if err := json.Unmarshal(b, &payload); err != nil {
		return noptions.NoVersion, false
	}
	return payload.Version, payload.Version != noptions.NoVersion
}
