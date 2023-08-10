package qfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moke-kit/mq/qiface"

	"moke-kit/mq/internal/nats"
)

type NatsResult struct {
	fx.Out
	qiface.MessageQueue `name:"NatsMQ"`
}

func (k *NatsResult) Execute(logger *zap.Logger, s SettingsParams) (err error) {
	k.MessageQueue, err = nats.NewMessageQueue(logger, s.NatsUrl)
	if err != nil {
		logger.Error("Nats message queue connect failure:",
			zap.Error(err),
			zap.String("address", s.NatsUrl))
	}
	return err
}

var NatsModule = fx.Provide(
	func(l *zap.Logger, s SettingsParams) (out NatsResult, err error) {
		err = out.Execute(l, s)
		return
	},
)

// For CLI testing purposes
func NewNatsMessageQueue(logger *zap.Logger, address string) (qiface.MessageQueue, error) {
	mq, err := nats.NewMessageQueue(logger, address)
	if err != nil {
		logger.Error("Nats message queue connect failure:",
			zap.Error(err),
			zap.String("address", address))
	}
	return mq, err
}
