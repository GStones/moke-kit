package nfx

import (
	"context"
	"fmt"
	"moke-kit/nsorm/nerrors"
	"net/url"

	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/nsorm/nosql/mongo"
)

type MongoParams struct {
	fx.In

	MongoClient *mongo2.Client `name:"MongoClient"`
}

type MongoResult struct {
	fx.Out

	MongoClient *mongo2.Client `name:"MongoClient"`
}

func (mr *MongoResult) NewDocument(
	lc fx.Lifecycle,
	l *zap.Logger,
	n SettingsParams,
) (err error) {
	if u, e := url.Parse(n.DocumentStoreUrl); e != nil {
		err = e
	} else if u.Scheme == "mongodb" {
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
		mr.MongoClient, err = mongo.NewMongoClient(conn, username, password)
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return mr.MongoClient.Disconnect(ctx)
			},
		})
	} else {
		return nerrors.ErrInvalidNosqlURL
	}
	return
}

var MongoPureModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		n SettingsParams,
	) (dOut MongoResult, err error) {
		err = dOut.NewDocument(lc, l, n)
		return
	},
)
