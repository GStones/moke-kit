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

func (k *LocalResult) init(logger *zap.Logger, s SettingsParams) error {
	k.Local = local.NewMessageQueue(logger, s.ChannelBufferSize, s.Persistent, s.BlockPublishUntilSubscriberAck)
	return nil
}

// CreateLocalModule creates a new local message queue module.
func CreateLocalModule(l *zap.Logger, s SettingsParams) (LocalResult, error) {
	out := LocalResult{}
	err := out.init(l, s)
	return out, err
}

// LocalModule is a module that provides the local message queue.
var LocalModule = fx.Provide(
	func(l *zap.Logger, s SettingsParams) (LocalResult, error) {
		return CreateLocalModule(l, s)
	},
)
