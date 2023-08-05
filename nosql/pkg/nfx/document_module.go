package nfx

import (
	"context"
	"net/url"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/document/mongo"
	"moke-kit/nosql/nerrors"
)

type DocumentStoreParams struct {
	fx.In
	DriverProvider diface.IDocumentProvider `name:"DriverProvider"`
}

type DocumentStoreResult struct {
	fx.Out
	DriverProvider diface.IDocumentProvider `name:"DriverProvider"`
}

func (g *DocumentStoreResult) NewDocument(
	lc fx.Lifecycle,
	l *zap.Logger,
	mParams MongoParams,
	n SettingsParams,
) (err error) {

	if mParams.MongoClient != nil {
		g.DriverProvider = mongo.NewProvider(mParams.MongoClient, l)
	}
	if n.DocumentStoreUrl != "" {
		if u, e := url.Parse(n.DocumentStoreUrl); e != nil {
			err = e
		} else {
			switch u.Scheme {
			case "mongodb":
				g.DriverProvider = mongo.NewProvider(mParams.MongoClient, l)
			case "test":
				l.Info("Connect to test", zap.String("url", "test"))
				//g.DriverProvider, err = mock.NewDocumentStoreProvider()

			default:
				return nerrors.ErrInvalidNosqlURL
			}
		}
	} else {
		return nerrors.ErrMissingNosqlURL
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
		mp MongoParams,
		sp SettingsParams,
	) (dOut DocumentStoreResult, err error) {
		err = dOut.NewDocument(lc, l, mp, sp)
		return
	},
)
