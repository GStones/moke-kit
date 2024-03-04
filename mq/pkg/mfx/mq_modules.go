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

func (g *MessageQueueResult) Execute(deploy utility.Deployments, i MQImplementations) (err error) {
	common.SetNamespace(deploy.String())
	g.MessageQueue = internal.NewMessageQueue(i.KafkaMQ, i.NatsMQ, i.NsqMQ, i.LocalMQ)

	return nil
}

var MqModule = fx.Provide(
	func(ap mfx.AppParams, i MQImplementations) (out MessageQueueResult, err error) {
		deployment := utility.ParseDeployments(ap.Deployment)
		err = out.Execute(deployment, i)
		return
	},
)
