package internal

import (
	document2 "github.com/gstones/platform/services/common/nosql/document"
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
	"io"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/json-iterator/go"

	"github.com/gstones/platform/services/common/dupblock"
	"github.com/gstones/platform/services/common/jsonx"
	"github.com/gstones/platform/services/common/jsonx/query"
	"go.uber.org/zap"
)

type document struct {
	version document2.Version
	data    []byte
	expires time.Time
}

func (doc *document) setTTL(ttl time.Duration) {
	if ttl > 0 {
		doc.expires = time.Now().Add(ttl)
	}
}

type DocumentStore struct {
	watcherManager document2.DocumentWatcherManager

	mutex     sync.RWMutex
	name      string
	owner     *DocumentStoreProvider
	documents map[document2.Key]*document
}

func (d *DocumentStore) Name() string {
	return d.name
}

func (d *DocumentStore) Contains(key document2.Key) (bool, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	_, ok := d.get(key)
	return ok, nil
}

func (d *DocumentStore) ListKeys(prefix string, opts ...document2.ScanOption) (list []document2.Key, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if o, e := document2.NewScanOptions(opts...); e != nil {
		err = e
		return
	} else {
		var offset int = o.Offset
		var limit int = o.Limit
		for k, _ := range d.documents {
			if strings.HasPrefix(k.String(), prefix) {
				if _, ok := d.get(k); ok {
					if offset > 0 {
						offset--
						continue
					}
					list = append(list, k)
					if limit > document2.ScanLimitNone && len(list) == limit {
						break
					}
				}
			}
		}
	}
	return
}

// Set creates or overwrites a document with the given keys and returns its version.  Use WithVersion to ensure this
// function is updating the version of the document that you expect.  If you don't use WithVersion then this
// function expects there to be no document.  If you want to set the document no matter what then use
// WithAnyVersion.
func (d *DocumentStore) Set(key document2.Key, opts ...document2.Option) (document2.Version, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if existingVersion, err := d.set(key, opts); err != nil {
		return document2.NoVersion, err
	} else {
		return existingVersion, nil
	}
}

// Internal version of Set(), for use only when the mutex is already locked, and there is no need to lock again
func (d *DocumentStore) set(key document2.Key, opts []document2.Option) (document2.Version, error) {
	existingDoc, docExists := d.get(key)
	if o, e := document2.NewOptions(opts...); e != nil {
		return document2.NoVersion, e
	} else {
		if data, e := jsonx.Marshal(o.Source); e != nil {
			return document2.NoVersion, e
		} else {
			if docExists {
				existingDoc.setTTL(o.TTL)
			}

			// TODO: rewrite this whole repetitive mess
			// Early out if WithAnyVersion is used - isn't really an early out
			if o.AnyVersion {
				if docExists {
					existingDoc.version++
					existingDoc.data = data
					return existingDoc.version, nil
				} else {
					d.documents[key] = &document{
						version: uint64(1),
						data:    data,
					}
					return uint64(1), nil
				}
			}

			if docExists {
				if o.Version == document2.NoVersion {
					return document2.NoVersion, errors2.ErrKeyExists
				} else if o.Version != existingDoc.version {
					return document2.NoVersion, errors2.ErrVersionMismatch
				} else {
					existingDoc.version++
					existingDoc.data = data
					return existingDoc.version, nil
				}
			} else {
				if o.Version == document2.NoVersion {
					d.documents[key] = &document{
						version: uint64(1),
						data:    data,
					}
					return uint64(1), nil
				} else {
					return document2.NoVersion, errors2.ErrKeyNotFound
				}
			}
		}
	}
}

// Get loads an existing document from the document store and returns its cas.  If no such document exists then
// this function fails.  Use WithTTL to update the document's expiration time.
func (d *DocumentStore) Get(key document2.Key, opts ...document2.Option) (document2.Version, error) {
	if o, e := document2.NewOptions(opts...); e != nil {
		return document2.NoVersion, e
	} else {
		d.mutex.RLock()
		defer d.mutex.RUnlock()
		if existingDoc, docExists := d.get(key); !docExists {
			return document2.NoVersion, errors2.ErrKeyNotFound
		} else {
			var err error

			existingDoc.setTTL(o.TTL)

			if o.Destination != nil {
				if err = jsonx.Unmarshal(existingDoc.data, o.Destination); err != nil {
					return document2.NoVersion, err
				}
			}
			return existingDoc.version, nil
		}
	}
}

func (d *DocumentStore) Scan(prefix string, destOpt document2.Option, scanOpts ...document2.ScanOption) (amt int, err error) {
	// make sure we have a destination
	if destOpt == nil {
		err = errors2.ErrDestIsNil
		return
	}

	opts := &document2.Options{}
	if e := destOpt(opts); e != nil {
		err = e
	} else if o, e := document2.NewScanOptions(scanOpts...); e != nil {
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

		d.mutex.Lock()
		defer d.mutex.Unlock()

		if res, e := d.scan(prefix, o.Query, o.Offset, o.Limit); e != nil {
			err = e
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
	}
	return
}

func (d *DocumentStore) scan(prefix string, queries []document2.ScanQuery, offset int, limit int) (res []interface{}, err error) {
	for k, v := range d.documents {
		// filter for prefix
		if !strings.HasPrefix(k.String(), prefix) {
			continue
		}

		// make sure we're not expired
		if _, ok := d.get(k); !ok {
			continue
		}

		buf := reflect.ValueOf(v.data).Bytes()
		good := true

		for _, q := range queries {
			switch q.ScanType {
			case document2.ScanTypeNOOP:
				continue
			case document2.ScanTypeKeyValue:
				good = jsoniter.Get(buf, q.KeyValue.Index).ToString() == q.KeyValue.Value
			case document2.ScanTypeKeyLike:
				good, err = regexp.MatchString(q.Regex, jsoniter.Get(buf, q.KeyValue.Index).ToString())
			case document2.ScanTypeRegex:
				if ok, e := regexp.Match(q.Regex, buf); e != nil {
					err = e
				} else if !ok {
					good = false
				}
			case document2.ScanTypeNum:
				vf := q.KeyValue.Value.(float64)
				qv := jsoniter.Get(buf, q.KeyValue.Index)
				if qv.ValueType() == jsoniter.NumberValue {
					qf := qv.ToFloat64()
					switch q.ScanOperation {
					case document2.ScanOpEquals:
						good = qf == vf
					case document2.ScanOpLessThan:
						good = qf < vf
					case document2.ScanOpGreaterThan:
						good = qf > vf
					default:
						err = errors.Wrap(errors2.ErrInternal, "unrecognized operation type")
					}
				} else {
					good = false
				}
			default:
				err = errors.Wrap(errors2.ErrNoScanType, "unrecognized scan type")
				break
			}

			if err != nil {
				return nil, err
			}
		}

		if good {
			// handle offset
			if offset > 0 {
				offset--
			} else {
				res = append(res, buf)
				if limit > document2.ScanLimitNone && len(res) >= limit {
					return
				}
			}
		}
	}

	return
}

// Internal implementation of Get that checks if a document should be expired before returning. If the document has an
// expired TTL set, then it will be deleted from the store. Mutex should already be locked before invocation.
func (d *DocumentStore) get(key document2.Key) (*document, bool) {
	if doc, ok := d.documents[key]; ok {
		if doc.expires.IsZero() || doc.expires.After(time.Now()) {
			return doc, true
		}
		// we're expired, delete
		delete(d.documents, key)
	}

	return nil, false
}

// Remove removes an existing document from the document store.  Use WithVersion to ensure this function is
// removing the version of the document that you expect.  Use WithAnyVersion to remove the document no matter what.
func (d *DocumentStore) Remove(key document2.Key, opts ...document2.Option) error {
	if o, e := document2.NewOptions(opts...); e != nil {
		return e
	} else {
		d.mutex.Lock()
		defer d.mutex.Unlock()

		if existingDoc, docExists := d.get(key); !docExists {
			return errors2.ErrKeyNotFound
		} else if o.AnyVersion {
			delete(d.documents, key)
			return nil
		} else if o.Version == existingDoc.version {
			delete(d.documents, key)
			return nil
		} else {
			return errors2.ErrVersionMismatch
		}
	}
}

// SetField sets a field (denoted by path) on a document and returns its cas.  If no such document exists then this
// function fails.  Use WithVersion to ensure this function is updating the version of the document that you expect.
func (d *DocumentStore) SetField(key document2.Key, path string, opts ...document2.Option) (document2.Version, error) {
	if o, e := document2.NewOptions(opts...); e != nil {
		return document2.NoVersion, e
	} else {
		var dest interface{}

		var resultVersion document2.Version
		if o.Version == document2.NoVersion {
			if o.AnyVersion {
				if beforeSet, err := d.Get(key, document2.WithAnyVersion(), document2.WithDestination(&dest)); err != nil {
					return document2.NoVersion, err
				} else if src, err := jsonx.Marshal(o.Source); err != nil {
					return document2.NoVersion, err
				} else if output, err := query.Execute(
					query.WithInput(dest),
					query.WithSetField([]byte(path), src),
				); err != nil {
					return document2.NoVersion, err
				} else {
					if resultVersion, err = d.Set(key,
						document2.WithVersion(beforeSet),
						document2.WithSource(output),
					); err != nil {
						return document2.NoVersion, err
					}
				}
			} else {
				return document2.NoVersion, errors.Wrap(errors2.ErrInternal, errors2.ErrInvalidVersioning.Error())
			}
		}
		return document2.Version(resultVersion), nil
	}
}

// AddWatcher adds a DocumentWatcher to the given keys.  Whenever changes to the document occur the watcher is told.
func (d *DocumentStore) AddWatcher(key document2.Key, watcher document2.DocumentWatcher) {
	if d.watcherManager.Logger == nil {
		d.watcherManager.Logger = zap.NewNop()
	}
	d.watcherManager.AddWatcher(key, watcher)
}

// RemoveWatcher removes a DocumentWatcher from the given keys.  It's the cleanup pair for AddWatcher.
func (d *DocumentStore) RemoveWatcher(key document2.Key, watcher document2.DocumentWatcher) {
	if d.watcherManager.Logger == nil {
		d.watcherManager.Logger = zap.NewNop()
	}
	d.watcherManager.RemoveWatcher(key, watcher)
}

// WatchersForKey returns a new slice of all DocumentWatcher instances for a keys.
func (d *DocumentStore) WatchersForKey(key document2.Key) []document2.DocumentWatcher {
	if d.watcherManager.Logger == nil {
		d.watcherManager.Logger = zap.NewNop()
	}
	return d.watcherManager.WatchersForKey(key)
}

// ApplyDUPBlock applies a document update block to the DocumentStore.
func (d *DocumentStore) ApplyDUPBlock(reader dupblock.Reader) error {
	var data []byte
	var outputResult interface{}
	var opts []query.Option
	var key document2.Key

	flush := func() error {
		if len(opts) > 0 {
			if version, err := d.Get(key,
				document2.WithAnyVersion(),
				document2.WithDestination(&data),
			); err != nil {
				return err
			} else {
				opts = append(opts, query.WithInput(data))

				if output, err := query.Execute(opts...); err != nil {
					return err
				} else if err := jsonx.Unmarshal(output, &outputResult); err != nil {
					return err
				} else if _, err := d.Set(key,
					document2.WithVersion(version),
					document2.WithSource(outputResult),
				); err != nil {
					return err
				} else {
					d.watcherManager.OnDocumentChanged(key, nil)
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
				} else if k, err := document2.NewKeyFromString(cmd.To); err != nil {
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
}

func (d *DocumentStore) PushBack(key document2.Key, path string, item interface{}) error {
	return errors2.ErrNotImplemented
}

// Touch touches a document, specifying a new expiry time for it.
func (d *DocumentStore) Touch(key document2.Key, opts ...document2.Option) (document2.Version, error) {
	if o, e := document2.NewOptions(opts...); e != nil {
		return document2.NoVersion, e
	} else {
		d.mutex.RLock()
		defer d.mutex.RUnlock()
		if existingDoc, docExists := d.get(key); !docExists {
			return document2.NoVersion, errors2.ErrKeyNotFound
		} else {
			existingDoc.setTTL(o.TTL)
			return existingDoc.version, nil
		}
	}
}

// Incr increments a numeric field by the provided amount.
func (d *DocumentStore) Incr(key document2.Key, path string, amount int32) error {
	return errors2.ErrNotImplemented
}
