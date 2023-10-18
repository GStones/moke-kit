package interceptors

import (
	"github.com/gstones/zinx/ziface"
	"go.uber.org/zap"
)

type LoggerInterceptor struct {
	logger *zap.Logger
}

func NewLoggerInterceptor(logger *zap.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{
		logger: logger,
	}
}

func (l *LoggerInterceptor) Intercept(chain ziface.IChain) ziface.IcResp {
	request := chain.Request()
	iRequest := request.(ziface.IRequest)
	l.logger.Info("zinx receive message", zap.Any("Data", iRequest.GetMessage()))
	return chain.Proceed(chain.Request())
}

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
