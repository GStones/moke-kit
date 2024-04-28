package interceptors

import (
	"context"

	"github.com/gstones/zinx/ziface"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/middlewares"
)

// RateLimitInterceptor rate limit interceptor
type RateLimitInterceptor struct {
	logger *zap.Logger
	rl     *middlewares.RateLimiter
}

// NewRateLimitInterceptor new rate limit interceptor
func NewRateLimitInterceptor(logger *zap.Logger, rateLimit int32) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		logger: logger,
		rl:     middlewares.CreateRateLimiter(int(rateLimit)),
	}
}

// Intercept intercept
func (r *RateLimitInterceptor) Intercept(chain ziface.IChain) ziface.IcResp {
	if err := r.rl.Limit(context.Background()); err != nil {
		r.logger.Error("rate limit exceeded", zap.Error(err))
	}
	return chain.Proceed(chain.Request())
}
