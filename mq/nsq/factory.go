package nsq

import (
	"github.com/gstones/platform/services/common/fx/fxsvcapp"
	"github.com/gstones/platform/services/common/mq"
	"github.com/gstones/platform/services/common/mq/nsq/internal"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Factory struct {
	fx.Out
	mq.MessageQueue `name:"NsqMQ"`
}

func (k *Factory) Execute(logger *zap.Logger, s fxsvcapp.GlobalSettings) (err error) {
	k.MessageQueue, err = internal.NewMessageQueue(logger, s.NsqConsumerUrl, s.NsqProducerUrl)
	if err != nil {
		logger.Error("Nsq message queue connect failure:",
			zap.Error(err),
			zap.String("consumer address", s.NsqConsumerUrl),
			zap.String("producer address", s.NsqProducerUrl))
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
func NewNsqMessageQueue(logger *zap.Logger, consumerAddress string, producerAddress string) (mq.MessageQueue, error) {
	mq, err := internal.NewMessageQueue(logger, consumerAddress, producerAddress)
	if err != nil {
		logger.Error("Nsq message queue connect failure:",
			zap.Error(err),
			zap.String("consumer address", consumerAddress),
			zap.String("producer address", producerAddress))
	}
	return mq, err
}
