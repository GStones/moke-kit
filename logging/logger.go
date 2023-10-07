package logging

import (
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogType string

const (
	LogTypeUndefined   LogType = ""
	LogTypeNone        LogType = "none"
	LogTypeDevelopment LogType = "dev"
	LogTypeProduction  LogType = "prod"
)

func ParseLogType(value string) LogType {
	switch LogType(value) {
	case LogTypeUndefined:
		return LogTypeDevelopment
	case LogTypeNone:
		return LogTypeNone
	case LogTypeDevelopment:
		return LogTypeDevelopment
	case LogTypeProduction:
		return LogTypeProduction
	default:
		panic(fmt.Errorf(`"%s" is an unknown log type`, value))
	}
}

func NewLogger(config Config) (logger *zap.Logger, err error) {
	switch ParseLogType(config.Type) {
	case LogTypeNone:
		logger = zap.NewNop()
	case LogTypeDevelopment:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = config.Build(zap.AddStacktrace(zap.FatalLevel))
	case LogTypeProduction:
		logger, err = zap.NewProduction(zap.AddStacktrace(zap.FatalLevel))
	}
	if logger != nil {
		logger.Info("log opened")
	}
	return
}

var Module = fx.Provide(
	LoadConfig,
	NewLogger,
)
