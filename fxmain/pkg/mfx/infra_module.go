package mfx

import (
	"go.uber.org/fx"
	"moke-kit/logging"
	nosql "moke-kit/nosql/pkg/module"
	server "moke-kit/server/pkg/module"
	tracing "moke-kit/tracing/module"
)

var InfraModule = fx.Options(
	AppModule,
	tracing.Module,
	server.Module,
	nosql.Module,
	logging.Module,
)
