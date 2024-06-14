package mfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

type SettingsParams struct {
	fx.In

	// Local channel buffer size
	// https://github.com/ThreeDotsLabs/watermill/blob/master/pubsub/gochannel/pubsub.go#L15
	ChannelBufferSize              int64 `name:"ChannelBufferSize"`
	Persistent                     bool  `name:"Persistent"`
	BlockPublishUntilSubscriberAck bool  `name:"BlockPublishUntilSubscriberAck"`

	// NatsUrl is the URL of the NATS server.
	NatsUrl string `name:"NatsUrl"`
}

type SettingsResult struct {
	fx.Out

	// Local Channel buffer size
	// https://github.com/ThreeDotsLabs/watermill/blob/master/pubsub/gochannel/pubsub.go#L15
	ChannelBufferSize              int64 `name:"ChannelBufferSize" envconfig:"CHANNEL_BUFFER_SIZE" default:"1024"`
	Persistent                     bool  `name:"Persistent" envconfig:"PERSISTENT" default:"false"`
	BlockPublishUntilSubscriberAck bool  `name:"BlockPublishUntilSubscriberAck" envconfig:"BLOCK_PUBLISH_UNTIL_SUBSCRIBER_ACK" default:"false"`

	// NatsUrl is the URL of the NATS server.
	NatsUrl string `name:"NatsUrl" envconfig:"NATS_URL" default:"nats://localhost:4222"`
}

func (ar *SettingsResult) loadFromEnv() (err error) {
	err = utility.Load(ar)
	return
}

// CreateSettingsModule creates a new settings module.
func CreateSettingsModule() (SettingsResult, error) {
	out := SettingsResult{}
	err := out.loadFromEnv()
	return out, err
}

// SettingModule is a module that provides the settings.
var SettingModule = fx.Provide(
	func() (out SettingsResult, err error) {
		return CreateSettingsModule()
	},
)
