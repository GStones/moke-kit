package module

import (
	"context"
	"testing"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
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
