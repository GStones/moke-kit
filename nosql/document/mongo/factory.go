package mongo

import (
	"context"
	"moke-kit/nosql/document/mongo/internal"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"moke-kit/nosql/document/diface"
)

func NewProvider(
	mClient *mongo.Client,
	logger *zap.Logger,
) diface.IDocumentProvider {
	return internal.NewDriverProvider(mClient, logger)
}

func NewMongoClient(
	connect string,
	username string,
	password string,
) (*mongo.Client, error) {
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
		return client, nil
	}
}
