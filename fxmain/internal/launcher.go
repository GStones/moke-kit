package internal

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/pkg/sfx"
)

type launchParams struct {
	fx.In
	ServiceBinder sfx.ServiceBinder

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

func Launch(in launchParams) (err error) {
	if err := in.ServiceBinder.Execute(in.Logger, in.Lifecycle); err != nil {
		return err
	}
	return err
}
