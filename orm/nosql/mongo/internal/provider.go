package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nosql/diface"
)

type DriverProvider struct {
	mClient *mongo.Client
	logger  *zap.Logger
}

func (dp *DriverProvider) Shutdown() error {
	return dp.mClient.Disconnect(context.Background())
}

func (dp *DriverProvider) OpenDbDriver(name string) (diface.ICollection, error) {
	db := dp.mClient.Database(name)
	if s, err := NewCollectionDriver(db); err != nil {
		return nil, err
	} else {
		return s, nil
	}
}

func NewDriverProvider(
	mClient *mongo.Client,
	logger *zap.Logger,
) *DriverProvider {
	return &DriverProvider{mClient, logger}
}
