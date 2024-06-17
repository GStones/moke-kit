package interceptors

import (
	"github.com/gstones/zinx/ziface"
	"go.uber.org/zap"
)

// LoggerInterceptor rate limit interceptor
type LoggerInterceptor struct {
	logger *zap.Logger
}

// NewLoggerInterceptor new rate limit interceptor
func NewLoggerInterceptor(logger *zap.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{
		logger: logger,
	}
}

// Intercept intercept
func (r *LoggerInterceptor) Intercept(chain ziface.IChain) ziface.IcResp {
	r.logger.Debug("zinx received message",
		zap.Any("msg", chain.GetIMessage()),
	)

	return chain.Proceed(chain.Request())
}
