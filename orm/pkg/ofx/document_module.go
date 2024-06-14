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

func (dsr *DocumentStoreResult) init(
	lc fx.Lifecycle,
	l *zap.Logger,
	mClient *mongo2.Client,
	connect string,
) error {
	if connect == "" {
		return nerrors.ErrMissingNosqlURL
	}

	if u, e := url.Parse(connect); e != nil {
		return e
	} else {
		switch u.Scheme {
		case "mongodb", "mongodb+srv":
			dsr.DriverProvider = mongo.NewProvider(mClient, l)
		default:
			return nerrors.ErrInvalidNosqlURL
		}
		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				return dsr.DriverProvider.Shutdown()
			},
		})
	}

	return nil
}

// CreateDocumentStore creates a new DocumentStoreResult.
func CreateDocumentStore(
	lc fx.Lifecycle,
	l *zap.Logger,
	mClient *mongo2.Client,
	connect string,
) (DocumentStoreResult, error) {
	var dOut DocumentStoreResult
	err := dOut.init(lc, l, mClient, connect)
	return dOut, err
}

// DocumentStoreModule provides  to the mfx dependency graph.
var DocumentStoreModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		mp MongoParams,
		sp SettingsParams,
	) (dOut DocumentStoreResult, err error) {
		return CreateDocumentStore(lc, l, mp.MongoClient, sp.DatabaseURL)
	},
)
