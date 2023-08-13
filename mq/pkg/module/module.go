package module

import (
	"go.uber.org/fx"

	"moke-kit/mq/pkg/qfx"
)

var Module = fx.Module("mq", fx.Options(
	qfx.MqModule,
	qfx.SettingModule,
))
