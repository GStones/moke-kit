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

// firebase auth middleware

// FirebaseAuthor is auth for grpc middleware
type FirebaseAuthor struct {
	client        *auth2.Client
	unAuthMethods map[string]struct{}
}

// Auth will auth every grpc request with firebase
// https://firebase.google.com/docs/auth/admin/verify-id-tokens
// NOTE: here we use firebase verifyIDToken to auth every grpc request,not to check the
// token has not been revoked or disabled.
// if you need to check the token has not been revoked please use VerifyIDTokenAndCheckRevoked to replace VerifyIDToken.
func (d *FirebaseAuthor) Auth(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	if _, ok := d.unAuthMethods[method]; ok {
		return context.WithValue(ctx, utility.WithOutTag, true), nil
	} else if token, err := auth.AuthFromMD(ctx, string(utility.TokenContextKey)); err != nil {
		return ctx, err
	} else if resp, err := d.client.VerifyIDToken(ctx, token); err != nil {
		return ctx, err
	} else {
		ctx = context.WithValue(ctx, utility.UIDContextKey, resp.UID)
		return ctx, nil
	}
}

// AddUnAuthMethod add unauth method
func (d *FirebaseAuthor) AddUnAuthMethod(method string) {
	if d.unAuthMethods == nil {
		d.unAuthMethods = make(map[string]struct{})
	}
	d.unAuthMethods[method] = struct{}{}
}

// FirebaseCheckModule is the firebase Auth module for grpc middleware
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
			client:        client,
			unAuthMethods: map[string]struct{}{},
		}
		return
	},
)
