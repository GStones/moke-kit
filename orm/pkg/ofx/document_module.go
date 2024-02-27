package ofx

import (
	"context"
	"net/url"

	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/mongo"
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
) (err error) {
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
			lc.Append(fx.Hook{
				OnStop: func(context.Context) error {
					return dsr.DriverProvider.Shutdown()
				},
			})
		}
	} else {
		return nerrors.ErrMissingNosqlURL
	}
	return
}

// DocumentStoreModule provides  to the mfx dependency graph.
var DocumentStoreModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		mp MongoParams,
		sp SettingsParams,
	) (dOut DocumentStoreResult, err error) {
		err = dOut.NewDocument(
			lc, l, mp.MongoClient,
			sp.DatabaseURL,
		)
		return
	},
)
