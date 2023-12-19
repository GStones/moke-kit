package logging

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/utility"
)

func NewLogger(deployment string) (logger *zap.Logger, err error) {
	if utility.ParseDeployments(deployment).IsProd() {
		if logger, err = zap.NewProduction(zap.AddStacktrace(zap.PanicLevel)); err != nil {
			return nil, err
		}
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		if logger, err = config.Build(zap.AddStacktrace(zap.ErrorLevel)); err != nil {
			return nil, err
		}
	}
	return
}

var Module = fx.Provide(
	func(params mfx.AppParams) (logger *zap.Logger, err error) {
		logger, err = NewLogger(params.Deployment)
		return
	},
)
