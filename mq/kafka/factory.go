package kafka

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/fx/fxsvcapp"
	"github.com/gstones/platform/services/common/mq"
	"github.com/gstones/platform/services/common/mq/kafka/internal"
)

type Factory struct {
	fx.Out
	mq.MessageQueue `name:"KafkaMQ"`
}

func (f *Factory) Execute(logger *zap.Logger, s fxsvcapp.GlobalSettings) (err error) {
	f.MessageQueue, err = internal.NewMessageQueue(logger, s.KafkaUrls)
	return err
}

var MQModule = fx.Provide(
	func(l *zap.Logger, s fxsvcapp.GlobalSettings) (out Factory, err error) {
		err = out.Execute(l, s)
		return
	},
)

// For CLI testing purposes
func NewKafkaMessageQueue(logger *zap.Logger, brokerUrls []string) (mq.MessageQueue, error) {
	return internal.NewMessageQueue(logger, brokerUrls)
}
