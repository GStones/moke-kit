package internal

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/pkg/module"
)

type launchParams struct {
	fx.In
	ServiceBinder module.ServiceBinder

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

// Launch launches the server
func Launch(in launchParams) (err error) {
	if err := in.ServiceBinder.Bind(in.Logger, in.Lifecycle); err != nil {
		return err
	}
	return err
}
