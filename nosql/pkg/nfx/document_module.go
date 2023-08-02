package nfx

import (
	"context"
	"fmt"
	"net/url"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/document/mongo"
	"moke-kit/nosql/errors"
)

type DocumentStoreParams struct {
	fx.In
	DriverProvider diface.IDocumentDb `name:"DriverProvider"`
}

type DocumentStoreResult struct {
	fx.Out
	DriverProvider diface.IDocumentDb `name:"DriverProvider"`
}

func (g *DocumentStoreResult) NewDocument(
	lc fx.Lifecycle,
	l *zap.Logger,
	n SettingsParams,
) (err error) {
	if n.DocumentStoreUrl != "" {
		if u, e := url.Parse(n.DocumentStoreUrl); e != nil {
			err = e
		} else {
			switch u.Scheme {
			case "mongodb":
				username := u.User.Username()
				if username == "" {
					username = n.NoSqlUser
				}

				password, set := u.User.Password()
				if !set {
					password = n.NoSqlPassword
				}
				conn := fmt.Sprintf("mongodb://%s", u.Host)
				l.Info("Connect to mongodb", zap.String("url", conn))
				g.DriverProvider, err = mongo.NewDriverProvider(conn, username, password, l)
			case "test":
				l.Info("Connect to test", zap.String("url", "test"))
				//g.DriverProvider, err = mock.NewDocumentStoreProvider()

			default:
				return errors.ErrInvalidNosqlURL
			}
		}
	} else {
		return errors.ErrMissingNosqlURL
	}

	if g.DriverProvider != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return g.DriverProvider.Shutdown()
			},
		})
	}
	return
}

var DocumentStoreModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		n SettingsParams,
	) (dOut DocumentStoreResult, err error) {
		err = dOut.NewDocument(lc, l, n)
		return
	},
)
