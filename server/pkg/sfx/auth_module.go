package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/server/siface"
)

type AuthServiceParams struct {
	fx.In

	AuthService siface.IAuth `name:"AuthService"`
}

type AuthServiceResult struct {
	fx.Out

	AuthService siface.IAuth `name:"AuthService"`
}
