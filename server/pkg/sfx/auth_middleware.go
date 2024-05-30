package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/server/siface"
)

// Auth module

type AuthMiddlewareParams struct {
	fx.In

	AuthMiddleware siface.IAuthMiddleware `name:"AuthMiddleware" optional:"true"`
}

type AuthMiddlewareResult struct {
	fx.Out

	AuthMiddleware siface.IAuthMiddleware `name:"AuthMiddleware"`
}
