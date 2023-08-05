package nfx

import (
	"context"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moke-kit/fxmain/pkg/mfx"
	"moke-kit/nosql/document/key"
	"net/url"

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
	mClient *mongo2.Client,
	connect string,
	deployment string,

) (err error) {
	key.SetNamespace(deployment)

	if mClient != nil {
		g.DriverProvider = mongo.NewProvider(mClient, l)
	}
	if connect != "" {
		if u, e := url.Parse(connect); e != nil {
			err = e
		} else {
			switch u.Scheme {
			case "mongodb":
				g.DriverProvider = mongo.NewProvider(mClient, l)
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
		as mfx.AppParams,
	) (dOut DocumentStoreResult, err error) {
		err = dOut.NewDocument(
			lc, l, mp.MongoClient,
			sp.DocumentStoreUrl, as.Deployment,
		)
		return
	},
)
