package interceptors

import (
	"github.com/gstones/zinx/ziface"
	"go.uber.org/zap"
)

type RecoverInterceptor struct {
	logger *zap.Logger
}

func NewRecoverInterceptor(logger *zap.Logger) *RecoverInterceptor {
	return &RecoverInterceptor{
		logger: logger,
	}
}

func (r *RecoverInterceptor) Intercept(chain ziface.IChain) ziface.IcResp {
	defer func() {
		if err := recover(); err != nil {
			r.logger.Error("panic error", zap.Any("err", err))
		}
	}()
	return chain.Proceed(chain.Request())
}
