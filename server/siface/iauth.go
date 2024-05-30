package siface

import "context"

// IAuthMiddleware is the interface for the auth middleware, which is used to authenticate the request.
type IAuthMiddleware interface {
	Auth(ctx context.Context) (context.Context, error)
	AddUnAuthMethod(method string)
}
