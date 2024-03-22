package mfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/internal/nats"
	"github.com/gstones/moke-kit/mq/miface"
)

type NatsResult struct {
	fx.Out
	NatsMQ miface.MessageQueue `name:"NatsMQ"`
}

func (k *NatsResult) Execute(logger *zap.Logger, s SettingsParams) (err error) {
	k.NatsMQ, err = nats.NewMessageQueue(logger, s.NatsUrl)
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
