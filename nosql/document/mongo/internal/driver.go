package internal

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/document/key"
	"moke-kit/nosql/nerrors"
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
		return diface.NoVersion, nerrors.ErrSourceIsNil
	} else {
		opt := options.Update()
		filter := bson.M{"_id": key.String()}
		if o.Version != diface.NoVersion {
			filter["version"] = o.Version
		} else {
			opt.SetUpsert(true)
		}
		update := bson.M{"$set": bson.M{"data": o.Source}, "$inc": bson.M{"version": 1}}
		if res, err := coll.UpdateOne(context.Background(), filter, update, opt); err != nil {
			return 0, err
		} else {
			if res.MatchedCount == 0 && o.Version != diface.NoVersion {
				return 0, nerrors.ErrVersionNotMatch
			} else {
				return o.Version + 1, nil
			}
		}
	}
}

func (c *DatabaseDriver) Get(key key.Key, opts ...diface.Option) (diface.Version, error) {
	coll := c.database.Collection(key.Prefix())
	if o, err := diface.NewOptions(opts...); err != nil {
		return diface.NoVersion, err
	} else {
		filter := bson.M{"_id": key.String()}
		if o.Version != diface.NoVersion {
			filter["version"] = o.Version
		}
		if res := coll.FindOne(context.Background(), filter); res.Err() != nil {
			if errors.Is(res.Err(), mongo.ErrNoDocuments) {
				return 0, nerrors.ErrNotFound
			}
			return 0, res.Err()
		} else {
			bRaw := &bson.Raw{}
			if err := res.Decode(bRaw); err != nil {
				return 0, err
			} else {
				if err := bRaw.Lookup("data").Unmarshal(o.Destination); err != nil {
					return 0, err
				}
				if err := bRaw.Lookup("version").Unmarshal(&o.Version); err != nil {
					return 0, err
				}
				return o.Version, nil
			}
		}
	}
}

func NewCollectionDriver(database *mongo.Database) (*DatabaseDriver, error) {
	return &DatabaseDriver{
		database: database,
	}, nil
}
