package ofx

import (
	"context"
	"fmt"
	"net/url"

	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/mongo"
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
	if n.DatabaseURL == "" {
		return nil
	}
	if u, e := url.Parse(n.DatabaseURL); e != nil {
		err = e
	} else if u.Scheme == "mongodb" {
		username := u.User.Username()
		if n.DatabaseUser != "" {
			username = n.DatabaseUser
		}
		password, _ := u.User.Password()
		if n.DatabasePassword != "" {
			password = n.DatabasePassword
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
		l.Error("Invalid mongodb url", zap.String("url", n.DatabaseURL))
		return nerrors.ErrInvalidNosqlURL
	}
	return
}

// MongoPureModule is the module for mongo driver
// https://github.com/mongodb/mongo-go-driver
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
