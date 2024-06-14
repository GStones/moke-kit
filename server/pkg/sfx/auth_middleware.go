package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/server/siface"
)

// Auth  middleware module

// AuthMiddlewareParams module params for injecting AuthMiddleware
type AuthMiddlewareParams struct {
	fx.In

	AuthMiddleware siface.IAuthMiddleware `name:"AuthMiddleware" optional:"true"`
}

// AuthMiddlewareResult module result for exporting AuthMiddleware
type AuthMiddlewareResult struct {
	fx.Out

	AuthMiddleware siface.IAuthMiddleware `name:"AuthMiddleware"`
}
