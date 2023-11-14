package interceptors

import (
	"github.com/gstones/zinx/ziface"
	"go.uber.org/zap"
)

// RateLimitInterceptor 流控拦截器
type RateLimitInterceptor struct {
	logger *zap.Logger
	rl     *rate.RateLimiter
}

func NewRateLimitInterceptor(logger *zap.Logger, rateLimit int32) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		logger: logger,
		rl:     rate.CreateRateLimiter(int(rateLimit)),
	}
}

// Intercept 拦截
func (r *RateLimitInterceptor) Intercept(chain ziface.IChain) ziface.IcResp {
	if err := r.rl.Limit(nil); err != nil {
		r.logger.Error("rate limit", zap.Error(err))
	}
	return chain.Proceed(chain.Request())
}
