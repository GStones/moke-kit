package module

import (
	"go.uber.org/fx"
	"moke-kit/fxmain/pkg/mfx"
	"moke-kit/logging"
	mq "moke-kit/mq/pkg/module"
	nosql "moke-kit/orm/pkg/module"
	server "moke-kit/server/pkg/module"
	tracing "moke-kit/tracing/module"
)

var AppModule = fx.Module("app",
	mfx.SettingModule,
	tracing.Module,
	server.Module,
	nosql.Module,
	logging.Module,
	mq.Module,
)
