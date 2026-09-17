package fxmain

import (
	"context"
	"time"

	_ "go.uber.org/automaxprocs" // Automatically set GOMAXPROCS:https://github.com/uber-go/automaxprocs

	"github.com/gstones/moke-kit/fxmain/internal"
	"github.com/gstones/moke-kit/fxmain/pkg/module"

	"go.uber.org/fx"
)

func appRun(base fx.Option, opts ...fx.Option) error {
	app := internal.NewApp(base, fx.Options(opts...))
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

// Main starts a batteries-included app (server + orm + mq + logging).
func Main(opts ...fx.Option) {
	if err := appRun(module.AppModule, opts...); err != nil {
		panic(err)
	}
}

// Core starts a thin app (settings + logging only). Pass server/orm/mq modules
// explicitly when those stacks are needed.
func Core(opts ...fx.Option) {
	if err := appRun(module.CoreModule, opts...); err != nil {
		panic(err)
	}
}
