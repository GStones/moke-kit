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

func (k *NatsResult) init(logger *zap.Logger, s SettingsParams) error {
	mq, err := nats.NewMessageQueue(logger, s.NatsUrl)
	if err != nil {
		logger.Error("Nats message queue connect failure:",
			zap.Error(err),
			zap.String("address", s.NatsUrl))
		return err
	}
	k.NatsMQ = mq
	return nil
}

// CreateNatsModule creates a new nats message queue module.
func CreateNatsModule(l *zap.Logger, s SettingsParams) (NatsResult, error) {
	out := NatsResult{}
	err := out.init(l, s)
	return out, err
}

// NatsModule is the module for nats message queue
var NatsModule = fx.Provide(
	func(l *zap.Logger, s SettingsParams) (NatsResult, error) {
		return CreateNatsModule(l, s)
	},
)
