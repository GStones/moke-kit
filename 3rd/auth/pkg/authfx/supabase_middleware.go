package authfx

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/gstones/moke-kit/server/pkg/sfx"
	"github.com/gstones/moke-kit/utility"
)

// SupabaseAuthor authenticates gRPC requests using Supabase.
// https://supabase.com/docs/guides/auth
type SupabaseAuthor struct {
	unauthTracker
	client *supabase.Client
}

// Auth authenticates every incoming gRPC request with Supabase.
func (d *SupabaseAuthor) Auth(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	if d.isUnauth(method) {
		return context.WithValue(ctx, utility.WithoutTag, true), nil
	} else if token, err := auth.AuthFromMD(ctx, string(utility.TokenContextKey)); err != nil {
		return ctx, err
	} else if resp, err := d.client.Auth.WithToken(token).GetUser(); err != nil {
		return ctx, err
	} else {
		ctx = context.WithValue(ctx, utility.UIDContextKey, resp.ID.String())
		return ctx, nil
	}
}

// SupabaseCheckModule is the Supabase auth module for the gRPC middleware.
var SupabaseCheckModule = fx.Provide(
	func(
		l *zap.Logger,
		sSetting SupabaseSettingParams,
	) (out sfx.AuthMiddlewareResult, err error) {
		c, err := supabase.NewClient(sSetting.URL, sSetting.Key, nil)
		if err != nil {
			return
		}
		out.AuthMiddleware = &SupabaseAuthor{
			unauthTracker: newUnauthTracker(),
			client:        c,
		}
		return
	},
)
