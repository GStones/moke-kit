package internal

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gstones/moke-kit/orm/nerrors"
	"github.com/gstones/moke-kit/orm/nosql/key"
	"github.com/gstones/moke-kit/orm/nosql/noptions"
)

// DatabaseDriver is a driver for a MongoDB database.
type DatabaseDriver struct {
	database *mongo.Database
}

// GetName Name returns the name of this ICollection.
func (dd *DatabaseDriver) GetName() string {
	return dd.database.Name()
}

// Set with a key and options
// If the version is not noVersion, then the version must match the version in the database,
// Or it will return error `ErrVersionNotMatch`
func (dd *DatabaseDriver) Set(ctx context.Context, key key.Key, opts ...noptions.Option) (noptions.Version, error) {
	o, err := noptions.NewOptions(opts...)
	if err != nil {
		return noptions.NoVersion, err
	}
	if o.Source == nil {
		return noptions.NoVersion, nerrors.ErrSourceIsNil
	}

	coll := dd.database.Collection(key.Prefix())
	filter := bson.M{"_id": key.String()}
	opt := options.Update()

	if o.Version != noptions.NoVersion {
		filter["version"] = o.Version
	} else {
		opt.SetUpsert(true)
	}

	update := bson.M{
		"$set": bson.M{"data": o.Source},
		"$inc": bson.M{"version": 1},
	}

	res, err := coll.UpdateOne(ctx, filter, update, opt)
	if err != nil {
		return 0, err
	}

	if res.MatchedCount == 0 && o.Version != noptions.NoVersion {
		return 0, nerrors.ErrVersionNotMatch
	}

	return o.Version + 1, nil
}

// Get data from mongoDB
func (dd *DatabaseDriver) Get(ctx context.Context, key key.Key, opts ...noptions.Option) (noptions.Version, error) {
	o, err := noptions.NewOptions(opts...)
	if err != nil {
		return noptions.NoVersion, err
	}

	coll := dd.database.Collection(key.Prefix())
	filter := bson.M{"_id": key.String()}
	if o.Version != noptions.NoVersion {
		filter["version"] = o.Version
	}

	var bRaw bson.Raw
	if err := coll.FindOne(ctx, filter).Decode(&bRaw); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, nerrors.ErrNotFound
		}
		return 0, err
	}

	if err := bRaw.Lookup("data").Unmarshal(o.Destination); err != nil {
		return 0, err
	}
	if err := bRaw.Lookup("version").Unmarshal(&o.Version); err != nil {
		return 0, err
	}

	return o.Version, nil
}

// Delete delete a document by a key
func (dd *DatabaseDriver) Delete(ctx context.Context, key key.Key) error {
	coll := dd.database.Collection(key.Prefix())
	filter := bson.M{"_id": key.String()}
	if _, err := coll.DeleteOne(ctx, filter); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nerrors.ErrNotFound
		}
		return err
	}
	return nil
}

// Incr increments a document from the nosql store. (tips: can not be used for document,because the version)
func (dd *DatabaseDriver) Incr(ctx context.Context, key key.Key, field string, amount int32) (int64, error) {
	coll := dd.database.Collection(key.Prefix())
	filter := bson.M{"_id": key.String()}
	update := bson.M{"$inc": bson.M{field: amount}}
	opt := options.FindOneAndUpdate()
	opt.SetUpsert(true)
	res := coll.FindOneAndUpdate(ctx, filter, update, opt)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return 0, nerrors.ErrNotFound
		}
		return 0, res.Err()
	}
	bRaw := &bson.Raw{}
	if err := res.Decode(bRaw); err != nil {
		return 0, err
	}
	var value int64
	if err := bRaw.Lookup(field).Unmarshal(&value); err != nil {
		return 0, err
	}
	return value, nil
}

// NewCollectionDriver creates a new DatabaseDriver.
func NewCollectionDriver(database *mongo.Database) (*DatabaseDriver, error) {
	return &DatabaseDriver{
		database: database,
	}, nil
}
