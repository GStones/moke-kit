package interceptors

import (
	"context"

	"github.com/gstones/zinx/ziface"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/common"
)

// RateLimitInterceptor 流控拦截器
type RateLimitInterceptor struct {
	logger *zap.Logger
	rl     *common.RateLimiter
}

func NewRateLimitInterceptor(logger *zap.Logger, rateLimit int32) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		logger: logger,
		rl:     common.CreateRateLimiter(int(rateLimit)),
	}
}

// Intercept 拦截
func (r *RateLimitInterceptor) Intercept(chain ziface.IChain) ziface.IcResp {
	if err := r.rl.Limit(context.Background()); err != nil {
		r.logger.Error("rate limit", zap.Error(err))
	}
	return chain.Proceed(chain.Request())
}
