package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/document/mongo/internal"
)

func NewDriverProvider(
	connect string,
	username string,
	password string,
	l *zap.Logger,
) (diface.IDocumentDb, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	credential := options.Credential{
		Username: username,
		Password: password,
	}
	cOptions := options.Client().ApplyURI(connect)
	if username != "" && password != "" {
		cOptions.SetAuth(credential)
	}
	if client, err := mongo.Connect(
		ctx,
		cOptions,
	); err != nil {
		return nil, err
	} else {
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			return nil, err
		}
		p := internal.NewProvider(client, l)
		return p, nil
	}
}
