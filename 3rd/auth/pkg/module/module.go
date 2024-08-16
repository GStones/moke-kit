package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/3rd/auth/pkg/authfx"
)

// SupabaseMiddlewareModule Provides supabase middleware for grpc
// if import this module, every grpc unary/stream will auth by supabase auth
// https://supabase.com/docs/guides/auth
var SupabaseMiddlewareModule = fx.Module("supabase_middleware",
	authfx.SupabaseSettingsModule,
	authfx.SupabaseCheckModule,
)

// FirebaseMiddlewareModule Provides firebase middleware for grpc
// if import this module, every grpc unary/stream will auth by firebase auth
// https://firebase.google.com/docs/auth/admin/verify-id-tokens
var FirebaseMiddlewareModule = fx.Module("firebase_middleware",
	authfx.FirebaseSettingsModule,
	authfx.FirebaseCheckModule,
)
