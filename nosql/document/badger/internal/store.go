package internal

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/dupblock"
	"github.com/gstones/platform/services/common/jsonx"
	"github.com/gstones/platform/services/common/jsonx/query"
	"github.com/gstones/platform/services/common/nosql/document"
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
)

type DocumentStore struct {
	document.DocumentWatcherManager

	name     string
	db       *badger.DB
	logger   *zap.Logger
	lockfile *os.File
	ticker   *time.Ticker
}

func openDB(opts badger.Options) (db *badger.DB, lockfile *os.File, err error) {
	lockfilename := filepath.Join(opts.Dir, "LOCK")

	// Ignore error, it's okay to fail, badger will report the LOCK file
	_ = os.Remove(lockfilename)
	if db, err = badger.Open(opts); err != nil {
		return nil, nil, err
	}

	return
}

func NewDocumentStore(dir string, name string, gcInterval time.Duration, logger *zap.Logger) (*DocumentStore, error) {
	d := path.Join(dir, name)
	opts := badger.DefaultOptions(d)

	// Per https://github.com/dgraph-io/badger/issues/476, must set Truncate true on Windows.
	if runtime.GOOS == "windows" {
		opts.Truncate = true
	}

	if err := os.MkdirAll(opts.Dir, os.ModePerm); err != nil {
		return nil, err
	} else if db, lockfile, err := openDB(opts); err != nil {
		return nil, convertBadgerError(err, "")
	} else {
		ds := &DocumentStore{
			name:     name,
			db:       db,
			lockfile: lockfile,
			logger:   logger,
			DocumentWatcherManager: document.DocumentWatcherManager{
				Logger: logger,
			},
		}
		ds.startGarbageCollection(gcInterval)
		return ds, nil
	}
}

func (d *DocumentStore) Name() string {
	return d.name
}

func (d *DocumentStore) Close() error {
	d.stopGarbageCollection()
	err := d.db.Close()

	d.lockfile.Close()
	d.stopGarbageCollection()
	return convertBadgerError(err, "")
}

// Contains checks to see if a document with the given keys exists.
func (d *DocumentStore) Contains(key document.Key) (ok bool, err error) {
	err = d.db.View(func(txn *badger.Txn) error {
		ok, err = d.contains(txn, key)
		return convertBadgerError(err, key.String())
	})
	return
}

func (d *DocumentStore) contains(txn *badger.Txn, key document.Key) (ok bool, err error) {
	if _, e := txn.Get(key.Bytes()); e != nil {
		if e == badger.ErrKeyNotFound {
			ok, err = false, nil
		} else {
			ok, err = false, e
		}
	} else {
		ok, err = true, nil
	}
	return
}

func (d *DocumentStore) ListKeys(prefix string, opts ...document.ScanOption) (list []document.Key, err error) {
	if o, e := document.NewScanOptions(opts...); e != nil {
		return nil, convertBadgerError(e, "")
	} else {
		var offset int = o.Offset
		var limit int = o.Limit

		if l, err := d.listKeys([]byte(prefix), offset, limit); err != nil {
			return nil, convertBadgerError(err, "")
		} else {
			return l, nil
		}
	}
}

func (d *DocumentStore) listKeys(prefix []byte, offset int, limit int) (list []document.Key, err error) {
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			if offset > 0 {
				offset--
				continue
			}
			item := it.Item()
			list = append(list, document.NewKeyFromBytesUnchecked(item.Key()))
			if limit > document.ScanLimitNone && len(list) >= limit {
				break
			}
		}
		return nil
	})
	return
}

// Set creates or overwrites  a document with the given keys and returns its version.  Use WithVersion to ensure this
// function is updating the version of the document that you expect.  If you don't use WithVersion then this
// function expects there to be no document.  If you want to set the document no matter what then use
// WithAnyVersion.
func (d *DocumentStore) Set(key document.Key, opts ...document.Option) (version document.Version, err error) {
	if o, e := document.NewOptions(opts...); e != nil {
		err = e
	} else if o.Source == nil {
		return document.NoVersion, errors2.ErrSourceIsNil
	} else if data, e := jsonx.Marshal(o.Source); e != nil {
		err = e
	} else if e := d.db.Update(func(txn *badger.Txn) error {
		version, err = d.set(txn, key, o.Version, o.AnyVersion, o.TTL, data)
		return err
	}); e != nil {
		version, err = document.NoVersion, convertBadgerError(e, key.String())
	} else {
		d.OnDocumentChanged(key, nil)
	}
	return
}

func (d *DocumentStore) set(
	txn *badger.Txn,
	key document.Key,
	expectedVersion document.Version,
	anyVersion bool,
	ttl time.Duration,
	data []byte,
) (version document.Version, err error) {
	if e := checkVersion(txn, key.Bytes(), expectedVersion, anyVersion); e != nil {
		err = e
	} else if e := setOrSetWithTTL(txn, key, ttl, data); e != nil {
		err = e
	} else if i, e := txn.Get(key.Bytes()); e != nil {
		err = e
	} else {
		version = document.Version(i.Version())
	}
	return
}

// Get loads an existing document from the document store and returns its cas.  If no such document exists then
// this function fails.  Use WithTTL to update the document's expiration time.
func (d *DocumentStore) Get(key document.Key, opts ...document.Option) (version document.Version, err error) {
	if o, e := document.NewOptions(opts...); e != nil {
		version, err = document.NoVersion, e
	} else if e := d.db.Update(func(txn *badger.Txn) error {
		version, err = d.get(txn, key, o.Version, o.TTL, o.Destination)
		return err
	}); e != nil {
		version, err = document.NoVersion, convertBadgerError(err, key.String())
	}
	return
}

func (d *DocumentStore) get(
	txn *badger.Txn,
	key document.Key,
	expectedVersion document.Version,
	ttl time.Duration,
	destination interface{},
) (version document.Version, err error) {
	if i, e := txn.Get(key.Bytes()); e != nil {
		return document.NoVersion, e
	} else if e := i.Value(func(val []byte) error {
		if ttl > 0 {
			if e := txn.Set(key.Bytes(), val); e != nil {
				return e
			}
		}
		return jsonx.Unmarshal(val, destination)
	}); e != nil {
		return document.NoVersion, e
	} else {
		return i.Version(), nil
	}
}

func (d *DocumentStore) Scan(prefix string, destOpt document.Option, scanOpts ...document.ScanOption) (amt int, err error) {
	// make sure we have a destination
	if destOpt == nil {
		err = errors2.ErrDestIsNil
		return
	}

	opts := &document.Options{}
	if e := destOpt(opts); e != nil {
		err = e
	} else if o, e := document.NewScanOptions(scanOpts...); e != nil {
		err = e
	} else if opts.Destination == nil {
		return 0, errors2.ErrDestIsNil
	} else if len(o.Query) == 0 {
		return 0, errors.Wrap(errors2.ErrNoScanType, "scan type unset")
	} else {
		v := reflect.ValueOf(opts.Destination).Elem()
		k := v.Kind()

		var dstLen int
		switch k {
		case reflect.Array:
			dstLen = v.Len()
		case reflect.Slice:
			dstLen = v.Cap()
		case reflect.Struct:
			dstLen = 1
		default:
			err = errors.Wrap(errors2.ErrInternal, "unsupported destination type")
			return
		}

		if o.Limit > dstLen {
			// TODO: log warning
			o.Limit = dstLen
		}

		if res, e := d.scan(prefix, o.Query, o.Offset, o.Limit); e != nil {
			err = convertBadgerError(err, "")
		} else if len(res) > 0 {
			for i, r := range res {
				raw := r.([]byte)
				if k == reflect.Struct {
					err = jsonx.Unmarshal(raw, opts.Destination)
					amt = 1
					return
				} else {
					t := v.Type().Elem()
					buf := reflect.New(t).Interface()
					if err = jsonx.Unmarshal(raw, &buf); err != nil {
						amt = 0
						return
					} else {
						elem := reflect.ValueOf(buf).Elem()
						v.Index(i).Set(elem)
						amt++
					}
				}
			}
		}
		return
	}
	return
}

func (d *DocumentStore) scan(prefix string, queries []document.ScanQuery, offset int, limit int) (res []interface{}, err error) {
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		pref := []byte(prefix)
		for it.Seek(pref); it.ValidForPrefix(pref); it.Next() {
			item := it.Item()
			if item.IsDeletedOrExpired() {
				continue
			}
			if err := item.Value(func(v []byte) error {
				buf := reflect.ValueOf(v).Bytes()
				good := true

				for _, q := range queries {
					switch q.ScanType {
					case document.ScanTypeNOOP:
						continue
					case document.ScanTypeKeyValue:
						if jsoniter.Get(buf, q.KeyValue.Index).ToString() != q.KeyValue.Value {
							good = false
						}
					case document.ScanTypeKeyLike:
						if ok, e := regexp.MatchString(q.Regex, jsoniter.Get(buf, q.KeyValue.Index).ToString()); e != nil {
							return e
						} else if !ok {
							good = false
						}
					case document.ScanTypeNum:
						vf := q.KeyValue.Value.(float64)
						qv := jsoniter.Get(buf, q.KeyValue.Index)
						if qv.ValueType() == jsoniter.NumberValue {
							qf := qv.ToFloat64()
							switch q.ScanOperation {
							case document.ScanOpEquals:
								good = qf == vf
							case document.ScanOpLessThan:
								good = qf < vf
							case document.ScanOpGreaterThan:
								good = qf > vf
							default:
								return errors.Wrap(errors2.ErrInternal, "unrecognized operation type")
							}
						} else {
							good = false
						}
					case document.ScanTypeRegex:
						if ok, e := regexp.Match(q.Regex, buf); e != nil {
							return e
						} else if !ok {
							good = false
						}
					default:
						return errors.Wrap(errors2.ErrNoScanType, "unrecognized scan type")
					}
				}
				if good {
					// handle offset
					if offset > 0 {
						offset--
					} else if vCopy, err := item.ValueCopy(make([]byte, len(v))); err != nil {
						return fmt.Errorf("unable to copy badger value: %w", err)
					} else {
						res = append(res, vCopy)
						if limit > document.ScanLimitNone && len(res) >= limit {
							return nil
						}
					}
				}
				return nil

			}); err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (d *DocumentStore) Remove(key document.Key, opts ...document.Option) error {
	if o, err := document.NewOptions(opts...); err != nil {
		return err
	} else if err := d.remove(key, o.Version, o.AnyVersion); err != nil {
		return convertBadgerError(err, key.String())
	} else {
		d.OnDocumentDeleted(key)
	}
	return nil
}

func (d *DocumentStore) remove(key document.Key, expectedVersion document.Version, anyVersion bool) error {
	return d.db.Update(func(txn *badger.Txn) error {
		if err := checkVersion(txn, key.Bytes(), expectedVersion, anyVersion); err != nil {
			return err
		} else if e := txn.Delete(key.Bytes()); e != nil {
			return e
		} else {
			return nil
		}
	})
}

func (d *DocumentStore) SetField(key document.Key, path string, opts ...document.Option) (document.Version, error) {
	if o, err := document.NewOptions(opts...); err != nil {
		return document.NoVersion, err
	} else if data, err := jsonx.Marshal(o.Source); err != nil {
		return document.NoVersion, err
	} else if version, err := d.setField(key, path, o.Version, o.AnyVersion, o.TTL, data); err != nil {
		return document.NoVersion, convertBadgerError(err, key.String())
	} else {
		d.OnDocumentChanged(key, nil)
		return version, nil
	}
}

func (d *DocumentStore) setField(
	key document.Key,
	path string,
	expectedVersion document.Version,
	anyVersion bool,
	ttl time.Duration,
	value []byte,
) (version document.Version, err error) {
	err = d.db.Update(func(txn *badger.Txn) error {
		var data []byte
		if version, err = d.get(txn, key, expectedVersion, ttl, &data); err != nil {
			return err
		} else if output, err := query.Execute(
			query.WithInput(data),
			query.WithSetField([]byte(path), value),
		); err != nil {
			return err
		} else if v, err := d.set(txn, key, version, anyVersion, ttl, output); err != nil {
			return err
		} else {
			version = v
			return nil
		}
	})
	return
}

func (d *DocumentStore) ApplyDUPBlock(reader dupblock.Reader) error {
	return d.db.Update(func(txn *badger.Txn) error {
		var data []byte
		var opts []query.Option
		var key document.Key

		flush := func() error {
			if len(opts) > 0 {
				if _, err := d.get(txn, key, document.NoVersion, 0, &data); err != nil {
					return convertBadgerError(err, key.String())
				} else {
					opts = append(opts, query.WithInput(data))
					if output, err := query.Execute(opts...); err != nil {
						return convertBadgerError(err, "")
					} else if _, err := d.set(txn, key, document.NoVersion, true, 0, output); err != nil {
						return convertBadgerError(err, key.String())
					} else {
						d.OnDocumentChanged(key, nil)
					}
				}

				opts = opts[:]
			}

			key.Clear()

			return nil
		}

		cmd := dupblock.Command{}

		for {
			if err := reader.Read(&cmd); err != nil {
				if err == io.EOF {
					break
				} else {
					return err
				}
			} else {
				switch cmd.Action {
				case dupblock.ActionSetKey:
					if err := flush(); err != nil {
						return err
					} else if k, err := document.NewKeyFromString(cmd.To); err != nil {
						return err
					} else {
						key = k
					}

				case dupblock.ActionSet:
					opts = append(opts, query.WithSetField([]byte(cmd.To), cmd.Value))
				case dupblock.ActionInsert:
					opts = append(opts, query.WithArrayInsert([]byte(cmd.To), cmd.Value))
				case dupblock.ActionIncrement:
					opts = append(opts, query.WithIncrement([]byte(cmd.To), cmd.Delta))
				case dupblock.ActionPushFront:
					opts = append(opts, query.WithArrayPushFront([]byte(cmd.To), cmd.Value))
				case dupblock.ActionPushBack:
					opts = append(opts, query.WithArrayPushBack([]byte(cmd.To), cmd.Value))
				case dupblock.ActionAddUnique:
					opts = append(opts, query.WithArrayUnique([]byte(cmd.To), cmd.Value))
				case dupblock.ActionDelete:
					opts = append(opts, query.WithDelete([]byte(cmd.To)))
				case dupblock.ActionCopy:
					opts = append(opts, query.WithCopy([]byte(cmd.From), []byte(cmd.To)))
				case dupblock.ActionMove:
					opts = append(opts, query.WithMove([]byte(cmd.From), []byte(cmd.To)))
				case dupblock.ActionSwap:
					opts = append(opts, query.WithSwap([]byte(cmd.From), []byte(cmd.To)))
				case dupblock.ActionUndefined:
					return dupblock.ErrUnknownDUPAction
				default:
					return dupblock.ErrUnhandledAction
				}
			}
		}

		return flush()
	})
}

func (d *DocumentStore) PushBack(key document.Key, path string, item interface{}) error {
	if e := d.db.Update(func(txn *badger.Txn) error {
		var data []byte
		var opts []query.Option

		if itemData, e := jsonx.Marshal(item); e != nil {
			return e
		} else {
			opts = append(opts, query.WithArrayPushBack([]byte(path), itemData))

			if _, err := d.get(txn, key, document.NoVersion, 0, &data); err != nil {
				return convertBadgerError(err, key.String())
			} else {
				opts = append(opts, query.WithInput(data))
				if output, err := query.Execute(opts...); err != nil {
					return convertBadgerError(err, "")
				} else if _, err := d.set(txn, key, document.NoVersion, true, 0, output); err != nil {
					return convertBadgerError(err, key.String())
				}
			}
		}

		return nil
	}); e != nil {
		return e
	}

	d.OnDocumentChanged(key, nil)
	return nil
}

func (d *DocumentStore) startGarbageCollection(c time.Duration) {
	d.ticker = time.NewTicker(c)
	go func() {
		for range d.ticker.C {
			for {
				if err := d.db.RunValueLogGC(0.5); err != nil {
					// Do not log error indicating nothing was garbage collected.
					if err != badger.ErrNoRewrite {
						d.logger.Info("Error during Badger garbage collection", zap.Error(err))
					}
					break
				}
			}
		}
	}()
}

func (d *DocumentStore) stopGarbageCollection() {
	d.ticker.Stop()
}

// Touch touches a document, specifying a new expiry time for it.
func (d *DocumentStore) Touch(key document.Key, opts ...document.Option) (version document.Version, err error) {
	if o, e := document.NewOptions(opts...); e != nil {
		version, err = document.NoVersion, e
	} else if e := d.db.Update(func(txn *badger.Txn) error {
		version, err = d.touch(txn, key, o.Version, o.TTL)
		return err
	}); e != nil {
		version, err = document.NoVersion, convertBadgerError(err, key.String())
	}
	return
}

func (d *DocumentStore) touch(
	txn *badger.Txn,
	key document.Key,
	expectedVersion document.Version,
	ttl time.Duration,
) (version document.Version, err error) {
	if i, e := txn.Get(key.Bytes()); e != nil {
		return document.NoVersion, e
	} else if e := i.Value(func(v []byte) error {
		if ttl > 0 {
			if e := txn.Set(key.Bytes(), v); e != nil {
				return e
			}
		}
		return nil
	}); e != nil {
		return document.NoVersion, e
	} else {
		return i.Version(), nil
	}
}

// Incr increments a numeric field by the provided amount.
func (d *DocumentStore) Incr(key document.Key, path string, amount int32) error {
	if e := d.db.Update(func(txn *badger.Txn) error {
		var data []byte
		var opts []query.Option

		opts = append(opts, query.WithIncrement([]byte(path), int(amount)))

		if _, err := d.get(txn, key, document.NoVersion, 0, &data); err != nil {
			return convertBadgerError(err, key.String())
		} else {
			opts = append(opts, query.WithInput(data))
			if output, err := query.Execute(opts...); err != nil {
				return convertBadgerError(err, "")
			} else if _, err := d.set(txn, key, document.NoVersion, true, 0, output); err != nil {
				return convertBadgerError(err, key.String())
			}
		}

		return nil
	}); e != nil {
		return e
	}

	d.OnDocumentChanged(key, nil)
	return nil
}
