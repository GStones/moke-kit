package nfx

import (
	"context"
	"net/url"

	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/fxmain/pkg/mfx"
	"moke-kit/gorm/nerrors"
	"moke-kit/gorm/nosql/diface"
	"moke-kit/gorm/nosql/key"
	"moke-kit/gorm/nosql/mongo"
)

type DocumentStoreParams struct {
	fx.In
	DriverProvider diface.IDocumentProvider `name:"DriverProvider"`
}

type DocumentStoreResult struct {
	fx.Out
	DriverProvider diface.IDocumentProvider `name:"DriverProvider"`
}

func (dsr *DocumentStoreResult) NewDocument(
	lc fx.Lifecycle,
	l *zap.Logger,
	mClient *mongo2.Client,
	connect string,
	deployment string,

) (err error) {
	key.SetNamespace(deployment)

	if mClient != nil {
		dsr.DriverProvider = mongo.NewProvider(mClient, l)
	}
	if connect != "" {
		if u, e := url.Parse(connect); e != nil {
			err = e
		} else {
			switch u.Scheme {
			case "mongodb":
				dsr.DriverProvider = mongo.NewProvider(mClient, l)
			case "test":
				l.Info("Connect to test", zap.String("url", "test"))
				//dsr.DriverProvider, err = mock.NewDocumentStoreProvider()

			default:
				return nerrors.ErrInvalidNosqlURL
			}
		}
	} else {
		return nerrors.ErrMissingNosqlURL
	}

	if dsr.DriverProvider != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return dsr.DriverProvider.Shutdown()
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
