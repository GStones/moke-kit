package cache

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/key"
)

func TestRedisCacheNilClientNoPanic(t *testing.T) {
	c := CreateRedisCache(zap.NewNop(), nil)
	k, err := key.NewKeyFromParts("a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.GetCache(context.Background(), k); got != nil {
		t.Fatal("expected cache miss")
	}
	if err := c.SetCache(context.Background(), k, map[string]any{"x": "y"}, time.Minute); err != nil {
		t.Fatal(err)
	}
	c.DeleteCache(context.Background(), k)
}
