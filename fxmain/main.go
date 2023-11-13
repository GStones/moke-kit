package fxmain

import (
	"context"
	"time"

	_ "go.uber.org/automaxprocs" // Automatically set GOMAXPROCS:https://github.com/uber-go/automaxprocs

	"github.com/gstones/moke-kit/fxmain/internal"
	"github.com/gstones/moke-kit/fxmain/pkg/module"

	"go.uber.org/fx"
)

func AppRun(opts ...fx.Option) error {
	app := internal.NewApp(opts...)
	if err := app.Run(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func Main(opts ...fx.Option) {
	if err := AppRun(
		module.AppModule,
		fx.Options(opts...),
		fx.Invoke(internal.Launch),
	); err != nil {
		panic(err)
	}
}
