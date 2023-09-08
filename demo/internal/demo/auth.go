package demo

import (
	"errors"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/pkg/sfx"
)

type DemoAuth struct {
}

func (d *DemoAuth) Auth(token string) (string, error) {
	if !strings.Contains(token, "test") {
		return "", errors.New("token error")
	}
	return token, nil
}

var AuthModule = fx.Provide(
	func(
		l *zap.Logger,
	) (out sfx.AuthServiceResult, err error) {
		out.AuthService = &DemoAuth{}
		return
	},
)
