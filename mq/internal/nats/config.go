package nats

import (
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
)

var jsConfig = nats.JetStreamConfig{Disabled: true}

var marshaler = &nats.JSONMarshaler{}
