package mfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/mq/internal/local"
	"github.com/gstones/moke-kit/mq/miface"
)

type LocalResult struct {
	fx.Out
	Local miface.MessageQueue `name:"LocalMQ"`
}

func (k *LocalResult) Execute(logger *zap.Logger, s SettingsParams) error {
	k.Local = local.NewMessageQueue(logger, s.ChannelBufferSize, s.Persistent, s.BlockPublishUntilSubscriberAck)
	return nil
}

var LocalModule = fx.Provide(
	func(l *zap.Logger, s SettingsParams) (out LocalResult, err error) {
		err = out.Execute(l, s)
		return
	},
)
