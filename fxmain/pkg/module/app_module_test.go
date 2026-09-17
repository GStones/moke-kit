package module

import (
	"context"
	"testing"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/utility"
)

func TestCoreModuleStartsWithoutServerOrORM(t *testing.T) {
	started := make(chan struct{})
	app := fx.New(
		CoreModule,
		fx.Invoke(func(logger *zap.Logger) {
			if logger == nil {
				t.Fatal("expected logger")
			}
			close(started)
		}),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Err(); err != nil {
		t.Fatalf("fx graph: %v", err)
	}
	if err := app.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal("invoke did not run")
	}
	if err := app.Stop(ctx); err != nil {
		t.Fatalf("stop: %v", err)
	}
}

func TestCoreAndAppShareDeploymentNamespace(t *testing.T) {
	t.Setenv("DEPLOYMENT", "PROD")
	t.Cleanup(func() { utility.SetNamespace("") })

	check := func(name string, opt fx.Option) {
		t.Run(name, func(t *testing.T) {
			var deploy string
			app := fx.New(
				fx.NopLogger,
				opt,
				fx.Invoke(func(p mfx.AppParams) {
					deploy = p.Deployment
				}),
			)
			if err := app.Err(); err != nil {
				t.Fatalf("fx graph: %v", err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := app.Start(ctx); err != nil {
				t.Fatalf("start: %v", err)
			}
			if deploy != "PROD" {
				t.Fatalf("Deployment = %q, want PROD", deploy)
			}
			if got := utility.Namespace(); got != "PROD" {
				t.Fatalf("Namespace() = %q, want PROD", got)
			}
			if err := app.Stop(ctx); err != nil {
				t.Fatalf("stop: %v", err)
			}
		})
	}

	check("core", CoreModule)
	check("app", AppModule)
}

func TestAppModuleConstructsWithNoServices(t *testing.T) {
	app := fx.New(AppModule)
	if err := app.Err(); err != nil {
		t.Fatalf("fx graph: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := app.Stop(ctx); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
