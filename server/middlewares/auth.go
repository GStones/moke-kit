package middlewares

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"

	"github.com/gstones/moke-kit/server/siface"
)

func authFunc(authClient siface.IAuthMiddleware) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if authClient != nil {
			return authClient.Auth(ctx)
		}
		return ctx, nil
	}
}

func allBut(_ context.Context, _ interceptors.CallMeta) bool {
	return true
}
