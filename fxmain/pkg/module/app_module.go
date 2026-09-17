package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/logging"
	mq "github.com/gstones/moke-kit/mq/pkg/module"
	nosql "github.com/gstones/moke-kit/orm/pkg/module"
	server "github.com/gstones/moke-kit/server/pkg/module"
)

// CoreModule is the thin app graph: settings + logging only.
// Compose server / orm / mq yourself (or use AppModule).
var CoreModule = fx.Module("core",
	mfx.SettingsModule,
	logging.Module,
)

// AppModule is the batteries-included graph used by fxmain.Main:
// settings, logging, server, orm, and mq settings/router.
var AppModule = fx.Module("app",
	CoreModule,
	server.Module,
	nosql.Module,
	mq.Module,
)
