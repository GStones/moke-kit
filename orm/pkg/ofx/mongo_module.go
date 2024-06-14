package ofx

import (
	"context"

	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"

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

func (mr *MongoResult) init(
	lc fx.Lifecycle,
	l *zap.Logger,
	n SettingsParams,
) error {
	if n.DatabaseURL == "" {
		return nil
	}
	cOptions := options.Client().ApplyURI(n.DatabaseURL)
	if n.DatabaseUser != "" {
		if cOptions.Auth == nil {
			cOptions.Auth = &options.Credential{}
		}
		cOptions.Auth.Username = n.DatabaseUser
	}
	if n.DatabasePassword != "" {
		if cOptions.Auth == nil {
			cOptions.Auth = &options.Credential{}
		}
		cOptions.Auth.Password = n.DatabasePassword
	}
	l.Info("Connecting mongodb", zap.String("url", n.DatabaseURL))
	client, err := mongo.NewMongoClient(cOptions)
	if err != nil {
		return err
	}
	mr.MongoClient = client
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return mr.MongoClient.Disconnect(ctx)
		},
	})

	return nil
}

// CreateMongoDriver creates a new Mongo driver client.
func CreateMongoDriver(
	lc fx.Lifecycle,
	l *zap.Logger,
	n SettingsParams,
) (MongoResult, error) {
	var dOut MongoResult
	err := dOut.init(lc, l, n)
	return dOut, err
}

// MongoPureModule is the module for mongo driver
// https://github.com/mongodb/mongo-go-driver
var MongoPureModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		n SettingsParams,
	) (MongoResult, error) {
		return CreateMongoDriver(lc, l, n)
	},
)
