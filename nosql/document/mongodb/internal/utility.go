package internal

import (
	"github.com/gstones/platform/services/common/nosql/document"
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/couchbase/gocbcore.v7"

	"github.com/snichols/gocb"
)

func getOrGetAndTouch(bucket *gocb.Bucket, key document.Key, expiry uint32, value interface{}) (version document.Version, err error) {
	if expiry == 0 {
		if v, e := bucket.Get(key.String(), value); e != nil {
			return document.NoVersion, e
		} else {
			return document.Version(v), nil
		}
	} else {
		if v, e := bucket.GetAndTouch(key.String(), expiry, value); e != nil {
			return document.NoVersion, e
		} else {
			return document.Version(v), nil
		}
	}
}

func cbExpiry(d time.Duration) uint32 {
	if d.Nanoseconds() == 0 {
		return 0
	} else {
		return uint32(d.Seconds())
	}
}

func convertMongodbError(e error, key string) error {
	if e == nil {
		return nil
	}

	switch e {
	case gocbcore.ErrKeyNotFound, mongo.ErrNoDocuments:
		if key == "" {
			return errors2.ErrKeyNotFound
		} else {
			return errors.Wrap(errors2.ErrKeyNotFound, key)
		}
	case gocbcore.ErrKeyExists:
		return errors2.ErrKeyExists
	default:
		return errors.Wrap(errors2.ErrDriverFailure, e.Error())
	}
}
