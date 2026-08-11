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
	var dst map[string]string
	if c.GetCache(context.Background(), k, &dst) {
		t.Fatal("expected cache miss")
	}
	c.SetCache(context.Background(), k, map[string]string{"x": "y"}, time.Minute)
	c.DeleteCache(context.Background(), k)
}
