package mock_kafka

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/fx/fxsvcapp"
	"github.com/gstones/platform/services/common/mq"
	"github.com/gstones/platform/services/common/mq/mock/internal"
)

type Factory struct {
	fx.Out
	mq.MessageQueue `name:"KafkaMQ"`
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
