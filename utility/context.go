package utility

import "context"

type ContextKey string

const (
	TokenContextKey ContextKey = "bearer"
	UIDContextKey   ContextKey = "uid"
)

func FromContext(ctx context.Context, key ContextKey) (string, bool) {
	value, ok := ctx.Value(key).(string)
	return value, ok
}

func NewContext(ctx context.Context, key ContextKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}
