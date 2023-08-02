package internal

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/document/key"
	"moke-kit/nosql/errors"
)

type DatabaseDriver struct {
	database *mongo.Database
}

func (c *DatabaseDriver) GetName() string {
	return c.database.Name()
}

func (c *DatabaseDriver) Set(key key.Key, opts ...diface.Option) (diface.Version, error) {
	coll := c.database.Collection(key.Prefix())
	if o, err := diface.NewOptions(opts...); err != nil {
		return diface.NoVersion, err
	} else if o.Source == nil {
		return diface.NoVersion, errors.ErrSourceIsNil
	} else if o.Version == diface.NoVersion {
		if _, err := coll.InsertOne(
			context.Background(),
			bson.M{"_id": key.String(), "version": 1, "source": o.Source}); err != nil {
			return 0, err
		} else {
			return 1, nil
		}
	} else if o.Version > 0 {
		filter := bson.M{"_id": key.String(), "version": o.Version}
		update := bson.M{"$set": o.Source, "$inc": bson.M{"version": 1}}
		if res, err := coll.UpdateOne(context.Background(), filter, update); err != nil {
			return 0, err
		} else {
			if res.MatchedCount == 0 {
				return 0, errors.ErrVersionNotMatch
			} else {
				return o.Version + 1, nil
			}
		}
	} else {
		return diface.NoVersion, errors.ErrSourceIsNil
	}
}

func (c *DatabaseDriver) Get(key key.Key, opts ...diface.Option) (diface.Version, error) {
	coll := c.database.Collection(key.Prefix())
	if o, err := diface.NewOptions(opts...); err != nil {
		return diface.NoVersion, err
	} else {
		filter := bson.M{"_id": key.String()}
		if o.Version > 0 {
			filter["version"] = o.Version
		}
		if err := coll.FindOne(context.Background(), filter).Decode(o.Destination); err != nil {
			return 0, err
		} else {
			return diface.NoVersion, nil
		}
	}
}

func NewCollectionDriver(database *mongo.Database) (*DatabaseDriver, error) {
	return &DatabaseDriver{
		database: database,
	}, nil
}
