package ofx

import (
	"errors"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestGormModuleRequiresDialector(t *testing.T) {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(zap.NewNop),
		GormModule,
		fx.Invoke(func(GormParams) {}),
	)
	if err := app.Err(); err == nil {
		t.Fatal("expected fx graph to fail without Dialector")
	}
}

func TestCreateGormDriverRejectsNilDialector(t *testing.T) {
	_, err := CreateGormDriver(nil, zap.NewNop(), nil)
	if !errors.Is(err, ErrGormDialectorRequired) {
		t.Fatalf("err = %v, want %v", err, ErrGormDialectorRequired)
	}
}
