package middlewares

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"

	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

func authFunc(authClient siface.IAuth) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		if token, err := auth.AuthFromMD(ctx, string(utility.TokenContextKey)); err != nil {
			return nil, err
		} else if authClient != nil {
			if uid, err := authClient.Auth(token); err != nil {
				return nil, err
			} else {
				return context.WithValue(ctx, utility.UIDContextKey, uid), nil
			}
		}
		return ctx, nil
	}
}

func allBut(_ context.Context, _ interceptors.CallMeta) bool {
	return true
}
