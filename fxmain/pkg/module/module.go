package module

import (
	"go.uber.org/fx"
	"moke-kit/fxmain/pkg/mfx"
)

var Module = fx.Module("app", fx.Options(
	mfx.AppModule,
))
