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
	defer func(chain ziface.IChain) {
		request := chain.Request()
		iRequest := request.(ziface.IRequest)
		l.logger.Info(
			"zinx receive",
			zap.String("Connection ID", iRequest.GetConnection().GetConnIdStr()),
			zap.String("IP", iRequest.GetConnection().RemoteAddr().String()),
			zap.Any("Data", iRequest.GetMessage()),
		)
	}(chain)

	return chain.Proceed(chain.Request())
}
