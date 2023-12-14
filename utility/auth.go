package utility

import (
	"context"
)

type Tag string

func (t Tag) String() string {
	return string(t)
}

const WithOutTag Tag = "auth.disabled"

// WithoutAuth overrides the default auth behavior and allows all methods to be called without an access token.
type WithoutAuth struct{}

// AuthFuncOverride allows all methods to be unauthenticated.
func (w *WithoutAuth) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	ctx = context.WithValue(ctx, WithOutTag, true)
	return ctx, nil
}
