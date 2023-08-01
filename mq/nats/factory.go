package nats

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/mq"
	"moke-kit/mq/nats/internal"
)

type Factory struct {
	fx.Out
	mq.MessageQueue `name:"NatsMQ"`
}

func (k *Factory) Execute(logger *zap.Logger, s fxsvcapp.GlobalSettings) (err error) {
	k.MessageQueue, err = internal.NewMessageQueue(logger, s.NatsUrl)
	if err != nil {
		logger.Error("Nats message queue connect failure:",
			zap.Error(err),
			zap.String("address", s.NatsUrl))
	}
	return err
}

var MQModule = fx.Provide(
	func(l *zap.Logger, s fxsvcapp.GlobalSettings) (out Factory, err error) {
		err = out.Execute(l, s)
		return
	},
)

// For CLI testing purposes
func NewNatsMessageQueue(logger *zap.Logger, address string) (mq.MessageQueue, error) {
	mq, err := internal.NewMessageQueue(logger, address)
	if err != nil {
		logger.Error("Nats message queue connect failure:",
			zap.Error(err),
			zap.String("address", address))
	}
	return mq, err
}
