package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/logging"
	mq "github.com/gstones/moke-kit/mq/pkg/module"
	nosql "github.com/gstones/moke-kit/orm/pkg/module"
	server "github.com/gstones/moke-kit/server/pkg/module"
)

// AppModule is the main module of the application.
// It includes some required modules like:
// SettingModule, server.Module, nosql.Module, logging.Module, and mq.Module.
var AppModule = fx.Module("app",
	mfx.SettingModule,
	server.Module,
	nosql.Module,
	logging.Module,
	mq.Module,
)
