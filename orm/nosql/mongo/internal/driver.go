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

// Set with a key and options.
// - No version and not AnyVersion: create only (fails if document exists).
// - WithVersion: CAS update (fails if version mismatches).
// - WithAnyVersion: overwrite/create regardless of current version.
func (dd *DatabaseDriver) Set(ctx context.Context, key key.Key, opts ...noptions.Option) (noptions.Version, error) {
	o, err := noptions.NewOptions(opts...)
	if err != nil {
		return noptions.NoVersion, err
	}
	if o.Source == nil {
		return noptions.NoVersion, nerrors.ErrSourceIsNil
	}

	coll := dd.database.Collection(key.Prefix())

	// Create-only: document must not already exist.
	if !o.AnyVersion && o.Version == noptions.NoVersion {
		doc := bson.M{
			"_id":     key.String(),
			"data":    o.Source,
			"version": int64(1),
		}
		if _, err := coll.InsertOne(ctx, doc); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return noptions.NoVersion, nerrors.ErrVersionNotMatch
			}
			return noptions.NoVersion, err
		}
		return 1, nil
	}

	filter := bson.M{"_id": key.String()}
	if !o.AnyVersion {
		filter["version"] = o.Version
	}

	update := bson.M{
		"$set": bson.M{"data": o.Source},
		"$inc": bson.M{"version": 1},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if o.AnyVersion {
		// $inc treats a missing field as 0, so upsert inserts version=1.
		opt.SetUpsert(true)
	}

	res := coll.FindOneAndUpdate(ctx, filter, update, opt)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return noptions.NoVersion, nerrors.ErrVersionNotMatch
		}
		return noptions.NoVersion, res.Err()
	}

	var out struct {
		Version noptions.Version `bson:"version"`
	}
	if err := res.Decode(&out); err != nil {
		return noptions.NoVersion, err
	}
	return out.Version, nil
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
	res, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return nerrors.ErrNotFound
	}
	return nil
}

// Incr increments a document from the nosql store. (tips: can not be used for document,because the version)
func (dd *DatabaseDriver) Incr(ctx context.Context, key key.Key, field string, amount int32) (int64, error) {
	coll := dd.database.Collection(key.Prefix())
	filter := bson.M{"_id": key.String()}
	update := bson.M{"$inc": bson.M{field: amount}}
	opt := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)
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
