package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/server/siface"
)

// Auth module

type AuthServiceParams struct {
	fx.In

	AuthService siface.IAuth `name:"AuthService" optional:"true"`
}

type AuthServiceResult struct {
	fx.Out

	AuthService siface.IAuth `name:"AuthService"`
}
