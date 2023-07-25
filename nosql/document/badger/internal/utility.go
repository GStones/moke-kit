package internal

import (
	"github.com/gstones/platform/services/common/nosql/document"
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
	"time"

	"github.com/pkg/errors"

	"github.com/dgraph-io/badger"
)

// Verifies that the Version value for the given keys matches the provided Version.  If the provided Version is not set
// then it's assumed to match.
func checkVersion(txn *badger.Txn, key []byte, version document.Version, anyVersion bool) error {
	if anyVersion {
		return nil
	} else if i, e := txn.Get(key); e != nil {
		if e == badger.ErrKeyNotFound {
			if version == document.NoVersion {
				return nil
			} else {
				return errors2.ErrKeyNotFound
			}
		} else {
			return e
		}
	} else if i.Version() != version {
		return errors2.ErrVersionMismatch
	} else {
		return nil
	}
}

// Sets the keys / value pair with an optional expiry time.
func setOrSetWithTTL(txn *badger.Txn, key document.Key, expiry time.Duration, data []byte) error {
	if expiry.Nanoseconds() == 0 {
		return txn.Set(key.Bytes(), data)
	} else {
		return txn.Set(key.Bytes(), data)
	}
}

func convertBadgerError(e error, key string) error {
	if e == nil {
		return nil
	}

	switch e {
	case badger.ErrKeyNotFound:
		if key == "" {
			return errors2.ErrKeyNotFound
		} else {
			return errors.Wrap(errors2.ErrKeyNotFound, key)
		}
	default:
		return errors.Wrap(e, errors2.ErrDriverFailure.Error())
	}
}
