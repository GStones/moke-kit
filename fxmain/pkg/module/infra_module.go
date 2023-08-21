package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/logging"
	mq "github.com/gstones/moke-kit/mq/pkg/module"
	nosql "github.com/gstones/moke-kit/orm/pkg/module"
	server "github.com/gstones/moke-kit/server/pkg/module"
	tracing "github.com/gstones/moke-kit/tracing/module"
)

var AppModule = fx.Module("app",
	mfx.SettingModule,
	tracing.Module,
	server.Module,
	nosql.Module,
	logging.Module,
	mq.Module,
)
