package authfx

import (
	"context"

	firebase "firebase.google.com/go/v4"
	auth2 "firebase.google.com/go/v4/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/server/pkg/sfx"
	"github.com/gstones/moke-kit/utility"
)

// FirebaseAuthor authenticates gRPC requests using Firebase ID tokens.
// https://firebase.google.com/docs/auth/admin/verify-id-tokens
//
// Note: VerifyIDToken is used here; it does not check for token revocation.
// Use VerifyIDTokenAndCheckRevoked if revocation checking is required.
type FirebaseAuthor struct {
	unauthTracker
	client *auth2.Client
}

// Auth authenticates every incoming gRPC request with Firebase.
func (d *FirebaseAuthor) Auth(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	if d.isUnauth(method) {
		return context.WithValue(ctx, utility.WithoutTag, true), nil
	} else if token, err := auth.AuthFromMD(ctx, string(utility.TokenContextKey)); err != nil {
		return ctx, err
	} else if resp, err := d.client.VerifyIDToken(ctx, token); err != nil {
		return ctx, err
	} else {
		ctx = context.WithValue(ctx, utility.UIDContextKey, resp.UID)
		return ctx, nil
	}
}

// FirebaseCheckModule is the Firebase auth module for the gRPC middleware.
var FirebaseCheckModule = fx.Provide(
	func(
		l *zap.Logger,
		sSetting FirebaseSettingParams,
	) (out sfx.AuthMiddlewareResult, err error) {
		ctx := context.Background()
		c, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(sSetting.KeyFile))
		if err != nil {
			return
		}
		client, err := c.Auth(ctx)
		if err != nil {
			return
		}
		out.AuthMiddleware = &FirebaseAuthor{
			unauthTracker: newUnauthTracker(),
			client:        client,
		}
		return
	},
)
