package mfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/mq/internal"
	"github.com/gstones/moke-kit/mq/miface"
)

type MessageQueueParams struct {
	fx.In

	MessageQueue miface.MessageQueue `name:"MessageQueue"`
}

type MessageQueueResult struct {
	fx.Out

	MessageQueue miface.MessageQueue `name:"MessageQueue"`
}

type MQImplementations struct {
	fx.In

	NatsMQ  miface.MessageQueue `name:"NatsMQ" optional:"true"`
	LocalMQ miface.MessageQueue `name:"LocalMQ" optional:"true"`
}

func (g *MessageQueueResult) init(mqs MQImplementations) (err error) {
	g.MessageQueue = internal.NewMessageQueue(mqs.NatsMQ, mqs.LocalMQ)
	return nil
}

// CreateMessageQueueModule creates a new message queue module.
// Namespace is owned by app settings; this does not overwrite it.
func CreateMessageQueueModule(mqs MQImplementations) (MessageQueueResult, error) {
	out := MessageQueueResult{}
	err := out.init(mqs)
	return out, err
}

// MqModule is a module that provides the message queue.
var MqModule = fx.Provide(
	func(mqs MQImplementations) (out MessageQueueResult, err error) {
		return CreateMessageQueueModule(mqs)
	},
)
