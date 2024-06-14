package mfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal"
	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/utility"
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
	KafkaMQ miface.MessageQueue `name:"KafkaMQ" optional:"true"`
	NsqMQ   miface.MessageQueue `name:"NsqMQ" optional:"true"`
	LocalMQ miface.MessageQueue `name:"LocalMQ" optional:"true"`
}

func (g *MessageQueueResult) init(mqs MQImplementations) (err error) {
	g.MessageQueue = internal.NewMessageQueue(mqs.KafkaMQ, mqs.NatsMQ, mqs.NsqMQ, mqs.LocalMQ)
	return nil
}

// CreateMessageQueueModule creates a new message queue module.
func CreateMessageQueueModule(deploy utility.Deployments, mqs MQImplementations) (MessageQueueResult, error) {
	common.SetNamespace(deploy.String())
	out := MessageQueueResult{}
	err := out.init(mqs)
	return out, err
}

// MqModule is a module that provides the message queue.
var MqModule = fx.Provide(
	func(ap mfx.AppParams, mqs MQImplementations) (out MessageQueueResult, err error) {
		deployment := utility.ParseDeployments(ap.Deployment)
		return CreateMessageQueueModule(deployment, mqs)
	},
)
