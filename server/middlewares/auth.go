package middlewares

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

// authFunc is a helper function to create a grpc auth interceptor
// that uses the provided authClient to authenticate incoming requests.
// In production deployments, a missing auth middleware fails closed.
func authFunc(authClient siface.IAuthMiddleware, deployments utility.Deployments) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if authClient != nil {
			return authClient.Auth(ctx)
		}
		if deployments.IsProd() {
			return nil, status.Error(codes.Unauthenticated, "auth middleware is required in production")
		}
		return ctx, nil
	}
}

func allBut(_ context.Context, _ interceptors.CallMeta) bool {
	return true
}
