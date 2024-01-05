package internal

import (
	"context"

	"go.uber.org/fx"
)

type App struct {
	*fx.App
}

func (app *App) Run() error {
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		app.Stop(ctx)
		return err
	}
	<-app.Done()
	return nil
}

func NewApp(opts ...fx.Option) *App {
	return &App{
		App: fx.New(
			fx.Options(opts...),
		),
	}
}
