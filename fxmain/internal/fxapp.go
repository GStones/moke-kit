package internal

import (
	"context"

	"go.uber.org/fx"
)

// App is a fx application
type App struct {
	*fx.App
}

// Run runs the application
func (app *App) Run() error {
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		app.Stop(ctx)
		return err
	}
	<-app.Done()
	return nil
}

// NewApp creates a new application
func NewApp(opts ...fx.Option) *App {
	return &App{
		App: fx.New(
			fx.Options(opts...),
		),
	}
}
