package siface

import "context"

type IAuthMiddleware interface {
	Auth(ctx context.Context) (context.Context, error)
	AddUnAuthMethod(method string)
}
