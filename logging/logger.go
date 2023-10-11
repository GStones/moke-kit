package logging

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/utility"
)

func NewLogger(deployment string) (logger *zap.Logger, err error) {
	switch utility.ParseDeployments(deployment) {
	case utility.DeploymentsLocal, utility.DeploymentsDev:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = config.Build(zap.AddStacktrace(zap.FatalLevel))
	case utility.DeploymentsProd:
		logger, err = zap.NewProduction(zap.AddStacktrace(zap.FatalLevel))
	}
	if logger != nil {
		logger.Info("log opened")
	}
	return
}

var Module = fx.Provide(
	func(params mfx.AppParams) (logger *zap.Logger, err error) {
		logger, err = NewLogger(params.Deployment)
		return
	},
)
