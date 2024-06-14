package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/diface"
	"github.com/gstones/moke-kit/orm/nosql/mongo/internal"
)

// NewProvider returns a new IDocumentProvider.
func NewProvider(
	mClient *mongo.Client,
	logger *zap.Logger,
) diface.IDocumentProvider {
	return internal.NewDriverProvider(mClient, logger)
}

// NewMongoClient returns a new mongo driver client.
func NewMongoClient(
	opts *options.ClientOptions,
) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if client, err := mongo.Connect(
		ctx,
		opts,
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
