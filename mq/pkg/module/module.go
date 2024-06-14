package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/mq/pkg/mfx"
)

// Module is the module for mq
var Module = fx.Module("mq",
	mfx.MqModule,
	mfx.SettingModule,
)
