package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"moke-kit/orm/nosql/diface"
)

type DriverProvider struct {
	mClient *mongo.Client
	logger  *zap.Logger
}

func (d *DriverProvider) Shutdown() error {
	return d.mClient.Disconnect(context.Background())
}

func (d *DriverProvider) OpenDbDriver(name string) (diface.ICollection, error) {
	db := d.mClient.Database(name)
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
