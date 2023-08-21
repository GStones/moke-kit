package qfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal"
	"github.com/gstones/moke-kit/mq/logic"
)

type MessageQueueParams struct {
	fx.In

	MessageQueue logic.MessageQueue `name:"MessageQueue"`
}

type MessageQueueResult struct {
	fx.Out

	MessageQueue logic.MessageQueue `name:"MessageQueue"`
}

type MQImplementations struct {
	fx.In

	NatsMQ  logic.MessageQueue `name:"NatsMQ" optional:"true"`
	KafkaMQ logic.MessageQueue `name:"KafkaMQ" optional:"true"`
	NsqMQ   logic.MessageQueue `name:"NsqMQ" optional:"true"`
	LocalMQ logic.MessageQueue `name:"LocalMQ" optional:"true"`
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

var MqModule = fx.Provide(
	func(s mfx.AppParams, i MQImplementations) (out MessageQueueResult, err error) {
		err = out.Execute(s, i)
		return
	},
)
