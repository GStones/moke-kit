package module

import (
	"context"
	"strings"
	"testing"

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/server/pkg/sfx"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestValidateSecurityConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  SecurityConfig
		wantErr string
	}{
		{
			name: "prod grpc requires auth",
			config: SecurityConfig{
				Deployment:      utility.ParseDeployments("prod"),
				HasGrpcServices: true,
			},
			wantErr: "auth middleware is required",
		},
		{
			name: "prod gateway requires auth",
			config: SecurityConfig{
				Deployment:       utility.ParseDeployments("prod_gcp"),
				HasGateway:       true,
				CorsAllowOrigins: "https://example.com",
			},
			wantErr: "auth middleware is required",
		},
		{
			name: "prod gateway requires explicit cors origins",
			config: SecurityConfig{
				Deployment: utility.ParseDeployments("prod"),
				HasGateway: true,
				HasAuth:    true,
			},
			wantErr: "CORS_ALLOW_ORIGINS",
		},
		{
			name: "prod gateway allows wildcard cors when explicit",
			config: SecurityConfig{
				Deployment:       utility.ParseDeployments("prod"),
				HasGateway:       true,
				HasAuth:          true,
				CorsAllowOrigins: "*",
			},
		},
		{
			name: "prod grpc with auth does not require cors",
			config: SecurityConfig{
				Deployment:      utility.ParseDeployments("prod"),
				HasGrpcServices: true,
				HasAuth:         true,
			},
		},
		{
			name: "prod without grpc or gateway does not require auth",
			config: SecurityConfig{
				Deployment: utility.ParseDeployments("prod"),
			},
		},
		{
			name: "dev grpc without auth is allowed",
			config: SecurityConfig{
				Deployment:      utility.ParseDeployments("dev"),
				HasGrpcServices: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSecurityConfig(tt.config)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("ValidateSecurityConfig() error = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("ValidateSecurityConfig() error = nil, want %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("ValidateSecurityConfig() error = %q, want containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestServiceBinderCheckSecurityConfigWarnsForNonProdMissingAuth(t *testing.T) {
	core, recorded := observer.New(zap.WarnLevel)
	logger := zap.New(core)
	binder := &ServiceBinder{
		AppParams: mfx.AppParams{
			Deployment: "dev",
		},
		GrpcServiceParams: sfx.GrpcServiceParams{
			GrpcServices: make([]siface.IGrpcService, 1),
		},
	}

	if err := binder.checkSecurityConfig(logger); err != nil {
		t.Fatalf("checkSecurityConfig() error = %v, want nil", err)
	}
	if recorded.Len() != 1 {
		t.Fatalf("recorded warnings = %d, want 1", recorded.Len())
	}
	if got := recorded.All()[0].Message; got != "auth middleware is not configured for grpc/gateway services" {
		t.Fatalf("warning message = %q", got)
	}
}

func TestServiceBinderCheckSecurityConfigFailsProdGatewayWithoutCors(t *testing.T) {
	binder := &ServiceBinder{
		AppParams: mfx.AppParams{
			Deployment: "prod-aws",
		},
		SettingsParams: sfx.SettingsParams{
			CorsAllowOrigins: " ",
		},
		AuthMiddlewareParams: sfx.AuthMiddlewareParams{
			AuthMiddleware: &testAuthMiddleware{},
		},
		GatewayServiceParams: sfx.GatewayServiceParams{
			GatewayServices: make([]siface.IGatewayService, 1),
		},
	}

	err := binder.checkSecurityConfig(zap.NewNop())
	if err == nil {
		t.Fatal("checkSecurityConfig() error = nil, want CORS error")
	}
	if !strings.Contains(err.Error(), "CORS_ALLOW_ORIGINS") {
		t.Fatalf("checkSecurityConfig() error = %q, want CORS_ALLOW_ORIGINS", err.Error())
	}
}

type testAuthMiddleware struct{}

func (*testAuthMiddleware) Auth(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (*testAuthMiddleware) AddUnAuthMethod(method string) {}
