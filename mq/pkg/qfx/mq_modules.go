package qfx

import (
	"go.uber.org/fx"
	"moke-kit/mq/common"
	"moke-kit/mq/qiface"

	"moke-kit/fxmain/pkg/mfx"
	"moke-kit/mq/internal"
)

type MessageQueueParams struct {
	fx.In

	MessageQueue qiface.MessageQueue `name:"MessageQueue"`
}

type MessageQueueResult struct {
	fx.Out

	MessageQueue qiface.MessageQueue `name:"MessageQueue"`
}

type MQImplementations struct {
	fx.In
	KafkaMQ qiface.MessageQueue `name:"KafkaMQ" optional:"true"`
	NatsMQ  qiface.MessageQueue `name:"NatsMQ" optional:"true"`
	NsqMQ   qiface.MessageQueue `name:"NsqMQ" optional:"true"`
	LocalMQ qiface.MessageQueue `name:"LocalMQ" optional:"true"`
}

func (g *MessageQueueResult) Execute(s mfx.AppParams, i MQImplementations) (err error) {
	common.SetNamespace(s.Deployment)

	// If run in TestMode, all Subscribe() and Publish() requests will be run through
	// the local:// mq implementation regardless of their chosen mq protocol
	if s.AppTestMode {
		g.MessageQueue = internal.NewMessageQueue(i.LocalMQ, i.LocalMQ, i.LocalMQ, i.LocalMQ)
	} else {
		g.MessageQueue = internal.NewMessageQueue(i.KafkaMQ, i.NatsMQ, i.NsqMQ, i.LocalMQ)
	}

	return nil
}

var Module = fx.Provide(
	func(s mfx.AppParams, i MQImplementations) (out MessageQueueResult, err error) {
		err = out.Execute(s, i)
		return
	},
)
