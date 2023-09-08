package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/mq/pkg/qfx"
)

var Module = fx.Module("mq",
	qfx.MqModule,
	qfx.SettingModule,
)
