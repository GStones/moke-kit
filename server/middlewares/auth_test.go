package middlewares

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gstones/moke-kit/utility"
)

func TestAuthFuncFailClosedInProd(t *testing.T) {
	fn := authFunc(nil, utility.DeploymentsProd)
	_, err := fn(context.Background())
	if err == nil {
		t.Fatal("expected error in prod without auth middleware")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got %v", status.Code(err))
	}
}

func TestAuthFuncAllowsNilInNonProd(t *testing.T) {
	fn := authFunc(nil, utility.DeploymentsLocal)
	ctx, err := fn(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx == nil {
		t.Fatal("expected context")
	}
}
