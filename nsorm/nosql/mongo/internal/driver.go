package internal

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"moke-kit/nsorm/nosql/noptions"

	"moke-kit/nsorm/nerrors"
	"moke-kit/nsorm/nosql/key"
)

type DatabaseDriver struct {
	database *mongo.Database
}

func (c *DatabaseDriver) GetName() string {
	return c.database.Name()
}

func (c *DatabaseDriver) Set(key key.Key, opts ...noptions.Option) (noptions.Version, error) {
	coll := c.database.Collection(key.Prefix())
	if o, err := noptions.NewOptions(opts...); err != nil {
		return noptions.NoVersion, err
	} else if o.Source == nil {
		return noptions.NoVersion, nerrors.ErrSourceIsNil
	} else {
		opt := options.Update()
		filter := bson.M{"_id": key.String()}
		if o.Version != noptions.NoVersion {
			filter["version"] = o.Version
		} else {
			opt.SetUpsert(true)
		}
		update := bson.M{"$set": bson.M{"data": o.Source}, "$inc": bson.M{"version": 1}}
		if res, err := coll.UpdateOne(context.Background(), filter, update, opt); err != nil {
			return 0, err
		} else {
			if res.MatchedCount == 0 && o.Version != noptions.NoVersion {
				return 0, nerrors.ErrVersionNotMatch
			} else {
				return o.Version + 1, nil
			}
		}
	}
}

func (c *DatabaseDriver) Get(key key.Key, opts ...noptions.Option) (noptions.Version, error) {
	coll := c.database.Collection(key.Prefix())
	if o, err := noptions.NewOptions(opts...); err != nil {
		return noptions.NoVersion, err
	} else {
		filter := bson.M{"_id": key.String()}
		if o.Version != noptions.NoVersion {
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
