package mock_nsq

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moke-kit/mq"
	"moke-kit/mq/mock/internal"
)

type Factory struct {
	fx.Out
	mq.MessageQueue `name:"NsqMQ"`
}

func (k *Factory) Execute(logger *zap.Logger, deployment string) (err error) {

	k.MessageQueue, err = internal.NewMessageQueue(logger, deployment)
	return err
}

var MQModule = fx.Provide(
	func(l *zap.Logger, s fxsvcapp.GlobalSettings) (out Factory, err error) {
		err = out.Execute(l, s.Deployment)
		return
	},
)

// For CLI testing purposes
func NewLocalMessageQueue(logger *zap.Logger, deployment string) (mq.MessageQueue, error) {
	return internal.NewMessageQueue(logger, deployment)
}
