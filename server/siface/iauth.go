package siface

import "context"

// IAuthMiddleware is the interface for the auth middleware, which is used to authenticate the request.
type IAuthMiddleware interface {
	// Auth is the method to authenticate the every grpc request
	Auth(ctx context.Context) (context.Context, error)
	// AddUnAuthMethod is the method to add the unauth method name
	AddUnAuthMethod(method string)
}
