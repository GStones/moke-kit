package internal

import (
	"fmt"
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/snichols/gocb"
	"go.uber.org/zap"
	"gopkg.in/couchbase/gocbcore.v7"

	"github.com/gstones/platform/services/common/dupblock"
	"github.com/gstones/platform/services/common/jsonx"
	"github.com/gstones/platform/services/common/nosql/document"
)

// DocumentStore adapts a Couchbase bucket to the nosql.DocumentStore interface.
type DocumentStore struct {
	mutex         sync.Mutex
	name          string
	owner         *DocumentStoreProvider
	dataBucket    *gocb.Bucket
	streamBucket  *gocb.StreamingBucket
	bucketWatcher *BucketWatcher
	logger        *zap.Logger
}

func NewDocumentStore(name string, owner *DocumentStoreProvider) (*DocumentStore, error) {
	return &DocumentStore{
		name:   name,
		owner:  owner,
		logger: owner.logger,
	}, nil
}

func (d *DocumentStore) Name() string {
	return d.name
}

// Contains checks to see if a document with the given keys exists in the DocumentStore.
func (d *DocumentStore) Contains(key document.Key) (contains bool, err error) {
	if b, e := d.data(); e != nil {
		err = convertCouchbaseError(e, key.String())
	} else {
		mo := gocbcore.GetMetaOptions{
			Key: key.Bytes(),
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		if _, e := b.IoRouter().GetMetaEx(mo, func(result *gocbcore.GetMetaResult, e error) {
			if e != nil && e == gocbcore.ErrKeyNotFound {
				contains = false
			} else if e != nil {
				err = e
			} else if result.Deleted != 0 {
				contains = false
			} else {
				contains = true
			}
			wg.Done()
		}); e != nil {
			wg.Done()
			err = convertCouchbaseError(e, key.String())
		}

		wg.Wait()
	}
	return
}

func (d *DocumentStore) ListKeys(prefix string, opts ...document.ScanOption) (list []document.Key, err error) {
	if b, e := d.data(); e != nil {
		return nil, e
	} else if o, e := document.NewScanOptions(opts...); e != nil {
		return nil, e
	} else {
		n1ql := fmt.Sprintf(
			"SELECT META(`%s`).id FROM `%s` WHERE META(`%s`).id LIKE $prefix",
			b.Name(), b.Name(), b.Name(),
		)
		params := map[string]interface{}{
			"prefix": fmt.Sprintf("%s%%", prefix),
		}
		if o.Offset > 0 {
			n1ql = fmt.Sprintf("%s OFFSET $offset", n1ql)
			params["offset"] = o.Offset
		}
		if o.Limit > 0 {
			n1ql = fmt.Sprintf("%s LIMIT $limit", n1ql)
			params["limit"] = o.Limit
		}

		if rows, e := d.performN1ql(n1ql, params, o.Timeout, false); e != nil {
			return nil, e
		} else {
			defer rows.Close()

			row := map[string]interface{}{}
			for rows.Next(&row) {
				// NB: if either of these errors ever fire, we're looking at pretty significant protocol level issues
				idRaw, ok := row["id"]
				if !ok {
					// we queried for one thing, and the rows that came back do not provide it
					return nil, errInternal(ErrN1QLResultIncomplete)
				}
				idStr, ok := idRaw.(string)
				if !ok {
					// it provided it, but we cannot coerce it to a string for some baffling reason
					return nil, errInternal(ErrN1QLResultFormatMismatch)
				}

				list = append(list, document.NewKeyFromStringUnchecked(idStr))
			}
			return list, nil
		}
	}
}

// Set creates or overwrites  a document with the given keys and returns its cas.  Use WithVersion to ensure this
// function is updating the version of the document that you expect.  If you don't use WithVersion then this
// function expects there to be no document.  If you want to set the document no matter what then use
// WithAnyVersion.
func (d *DocumentStore) Set(key document.Key, opts ...document.Option) (document.Version, error) {
	if o, e := document.NewOptions(opts...); e != nil {
		return document.NoVersion, e
	} else if o.Source == nil {
		return document.NoVersion, errors2.ErrSourceIsNil
	} else if b, e := d.data(); e != nil {
		return document.NoVersion, convertCouchbaseError(e, "")
	} else {
		if o.Version == document.NoVersion {
			if o.AnyVersion {
				if version, e := b.Upsert(key.String(), o.Source, cbExpiry(o.TTL)); e != nil {
					return document.NoVersion, convertCouchbaseError(e, key.String())
				} else {
					return document.Version(version), nil
				}
			} else {
				if version, e := b.Insert(key.String(), o.Source, cbExpiry(o.TTL)); e != nil {
					if e == gocb.ErrKeyExists {
						return document.NoVersion, errors2.ErrVersionMismatch
					} else {
						return document.NoVersion, convertCouchbaseError(e, key.String())
					}
				} else {
					return document.Version(version), nil
				}
			}
		} else if version, e := b.Replace(key.String(), o.Source, gocb.Cas(o.Version), cbExpiry(o.TTL)); e != nil {
			if e == gocbcore.ErrInvalidArgs {
				return document.NoVersion, errors2.ErrVersionMismatch
			} else {
				return document.NoVersion, convertCouchbaseError(e, key.String())
			}
		} else {
			return document.Version(version), nil
		}
	}
}

// Get loads an existing document from the document store and returns its cas.  If no such document exists then
// this function fails.
func (d *DocumentStore) Get(key document.Key, opts ...document.Option) (document.Version, error) {
	if o, e := document.NewOptions(opts...); e != nil {
		return document.NoVersion, e
	} else if b, e := d.data(); e != nil {
		return document.NoVersion, convertCouchbaseError(e, "")
	} else {
		if version, e := getOrGetAndTouch(b, key, cbExpiry(o.TTL), o.Destination); e != nil {
			return document.NoVersion, convertCouchbaseError(e, key.String())
		} else {
			return document.Version(version), nil
		}
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

		if res, e := d.scan(prefix, o.Query, o.Offset, o.Limit, o.Timeout); e != nil {
			err = convertCouchbaseError(e, "")
		} else if len(res) > 0 {
			amt = len(res)
			for i, r := range res {
				raw := r.([]byte)
				if k == reflect.Struct {
					err = jsonx.Unmarshal(raw, opts.Destination)
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
					}
				}
			}
		}
		return
	}
	return
}

func (d *DocumentStore) scan(prefix string, queries []document.ScanQuery, offset int, limit int, timeout time.Duration) (res []interface{}, err error) {
	if b, e := d.data(); e != nil {
		return nil, e
	} else {
		params := map[string]interface{}{
			"prefix": fmt.Sprintf("%s%%", prefix),
		}
		n1ql := fmt.Sprintf(
			"SELECT * FROM `%s` WHERE META(`%s`).id LIKE $prefix",
			b.Name(), b.Name(),
		)

		for n, q := range queries {
			var op string
			switch q.ScanType {
			case document.ScanTypeNOOP:
				continue
			case document.ScanTypeKeyValue:
				op = "="
			case document.ScanTypeKeyLike:
				op = "LIKE"
			case document.ScanTypeNum:
				switch q.ScanOperation {
				case document.ScanOpEquals:
					op = "="
				case document.ScanOpLessThan:
					op = "<"
				case document.ScanOpGreaterThan:
					op = ">"
				default:
					return nil, errors.Wrap(errors2.ErrInternal, "unrecognized scan operation")
				}
			case document.ScanTypeRegex:
				return nil, errors.Wrap(errors2.ErrNotImplemented, "unsupported scan type")
			default:
				return nil, errors.Wrap(errors2.ErrInternal, "unsupported scan type")
			}

			n1ql = fmt.Sprintf("%s AND `%s` %s $val%d", n1ql, q.KeyValue.Index, op, n)
			params[fmt.Sprint("val", n)] = q.KeyValue.Value
		}

		if offset > 0 {
			n1ql = fmt.Sprintf("%s OFFSET $offset", n1ql)
			params["offset"] = offset
		}
		if limit > 0 {
			n1ql = fmt.Sprintf("%s LIMIT $limit", n1ql)
			params["limit"] = limit
		}

		if rows, e := d.performN1ql(n1ql, params, timeout, false); e != nil {
			return nil, e
		} else {
			defer rows.Close()

			// TODO: Rewrite the entire return path here...
			/**
			 * Yes, we're putting this BACK into json for now.
			 *
			 * The problem is that unless we have the exact target struct, gocb returns the data as a map (I cannot get
			 * it to correctly populate a reflection-instantiated struct instead). Our primary alternative to passing
			 * the data through json here seems to be using hashicorp's mapstructure lib instead. - al[2019.02.14]
			 */
			row := map[string]interface{}{}
			for rows.Next(&row) {
				if obj, e := jsonx.Marshal(row[b.Name()]); e != nil {
					return nil, e
				} else {
					res = append(res, obj)
				}
			}
		}
	}
	return
}

/**
 * NB: This method is way incomplete - and is unused.
 *     Final solution is blocked on #1148, leaving this in but commented out until then.
 *
func (d *DocumentStore) performRegex(indexName string, search string, timeout time.Duration) (gocb.SearchResults, error) {
	if b, e := d.data(); e != nil {
		return nil, e
	} else {
		query := gocb.NewSearchQuery(indexName, cbft.NewRegexpQuery(search))
		// check if the context specifies a timeout
		if timeout != nosql.ScanTimeoutNone {
			query.Timeout(timeout)
		}

		if rows, e := b.ExecuteSearchQuery(query); e != nil {
			return nil, e
		} else {
			return rows, nil
		}
	}

}
*/

// Consolidates some of the boilerplate of performing an n1ql query. NB: As we are not processing the result here, the
// invoking method should close the result object when they are done with it instead.
func (d *DocumentStore) performN1ql(n1ql string, params interface{}, timeout time.Duration, adhoc bool) (gocb.QueryResults, error) {
	if b, e := d.data(); e != nil {
		return nil, e
	} else {
		query := gocb.NewN1qlQuery(n1ql)
		// adhoc queries are not prepared/cached before execution
		query.AdHoc(adhoc)
		// check if the context specifies a timeout
		if timeout != document.ScanTimeoutNone {
			query.Timeout(timeout)
		}

		if rows, e := b.ExecuteN1qlQuery(query, params); e != nil {
			return nil, e
		} else {
			return rows, nil
		}
	}
}

// Remove removes an existing document from the document store.  Use WithVersion to ensure this function is
// removing the version of the document that you expect.  Use WithAnyVersion to remove the document no matter what.
func (d *DocumentStore) Remove(key document.Key, opts ...document.Option) error {
	if o, e := document.NewOptions(opts...); e != nil {
		return e
	} else if b, e := d.data(); e != nil {
		return convertCouchbaseError(e, "")
	} else {
		cas := gocb.Cas(o.Version)
		if o.AnyVersion {
			cas = 0
		}
		if _, e := b.Remove(key.String(), cas); e != nil {
			return convertCouchbaseError(e, key.String())
		} else {
			return nil
		}
	}
}

func (d *DocumentStore) SetField(key document.Key, path string, opts ...document.Option) (document.Version, error) {
	if o, e := document.NewOptions(opts...); e != nil {
		return document.NoVersion, e
	} else if b, e := d.data(); e != nil {
		return document.NoVersion, convertCouchbaseError(e, "")
	} else {
		var builder *gocb.MutateInBuilder
		if o.Version == document.NoVersion {
			if o.AnyVersion {
				builder = b.MutateIn(key.String(), 0, 0)
			} else {
				return document.NoVersion, errInternal(errors2.ErrInvalidVersioning)
			}
		} else {
			builder = b.MutateIn(key.String(), gocb.Cas(o.Version), 0)
		}

		builder = builder.Upsert(
			path,
			o.Source,
			true,
		)
		if f, e := builder.Execute(); e != nil {
			return document.NoVersion, convertCouchbaseError(e, "")
		} else {
			return document.Version(f.Cas()), nil
		}
	}
}

func (d *DocumentStore) ApplyDUPBlock(reader dupblock.Reader) error {
	if b, e := d.data(); e != nil {
		return e
	} else {
		var cmd dupblock.Command
		var cmds []dupblock.Command
		var key document.Key

		// Would be nice if we could guarantee ActionSetKey was first.
		keySet := false
		for {
			if err := reader.Read(&cmd); err != nil {
				if err != io.EOF {
					return err
				}
				break
			} else {
				cmds = append(cmds, cmd)
				if !keySet && cmd.Action == dupblock.ActionSetKey {
					if key, err = document.NewKeyFromString(cmd.To); err != nil {
						return err
					} else {
						keySet = true
					}
				}
			}
		}
		builder := b.MutateIn(key.String(), 0, 0)
		cmdAdded := false

		for _, cmd := range cmds {
			switch cmd.Action {
			case dupblock.ActionUndefined:
				return errors2.ErrUnknownDUPAction
			case dupblock.ActionSet:
				builder = builder.Upsert(
					cmd.To,
					cmd.Value,
					true,
				)
				cmdAdded = true
			case dupblock.ActionInsert:
				builder = builder.ArrayInsert(
					cmd.To,
					cmd.Value,
				)
				cmdAdded = true
			case dupblock.ActionIncrement:
				var num int
				if f, e := b.LookupIn(key.String()).Get(cmd.To).Execute(); e != nil {
					return e
				} else if e = f.Content(cmd.To, &num); e != nil {
					return e
				} else {
					num += cmd.Delta
					builder = builder.Upsert(
						cmd.To,
						num,
						true,
					)
					cmdAdded = true
				}
			case dupblock.ActionPushFront:
				builder = builder.ArrayPrepend(
					cmd.To,
					cmd.Value,
					true,
				)
				cmdAdded = true
			case dupblock.ActionPushBack:
				builder = builder.ArrayAppend(
					cmd.To,
					cmd.Value,
					true,
				)
				cmdAdded = true
			case dupblock.ActionAddUnique:
				// From https://developer.couchbase.com/documentation/server/5.5/sdk/subdocument-operations.html
				// Note that currently the addunique will fail with a Path Mismatch error if the array contains JSON floats, objects, or arrays.
				// The addunique operation will also fail with Cannot Insert if the value to be added is one of those types as well.
				var x interface{}
				if f, e := b.LookupIn(key.String()).Get(cmd.To).Execute(); e != nil {
					return e
				} else if e = f.Content(cmd.To, &x); e != nil {
					return e
				} else {
					switch p := x.(type) {
					case []interface{}:
						insert := true
						var value map[string]interface{}
						if err := jsonx.Unmarshal(cmd.Value, &value); err != nil {
							return errInternal(ErrUnknownValueType)
						}
						for _, a := range p {
							if !reflect.DeepEqual(a, value) {
								insert = false
								break
							}
						}
						if insert {
							builder = builder.ArrayPrepend(
								cmd.To,
								cmd.Value,
								true,
							)
							cmdAdded = true
						}
					default:
						return errInternal(ErrInvalidDestination)
					}
				}
			case dupblock.ActionDelete:
				builder = builder.Remove(cmd.To)
				cmdAdded = true
			case dupblock.ActionCopy:
				var x interface{}
				if f, e := b.LookupIn(key.String()).Get(cmd.From).Execute(); e != nil {
					return e
				} else if e = f.Content(cmd.From, &x); e != nil {
					return e
				} else {
					// What datatype is destination? Arrays are handled differently in Couchbase. Yay!
					if cmd.To[len(cmd.To)-1] == ']' {
						builder = builder.ArrayInsert(cmd.To, x)
					} else {
						builder = builder.Upsert(cmd.To, x, true)
					}
					cmdAdded = true
				}
			case dupblock.ActionMove:
				var x interface{}
				if f, e := b.LookupIn(key.String()).Get(cmd.From).Execute(); e != nil {
					return e
				} else if e = f.Content(cmd.From, &x); e != nil {
					return e
				} else {
					// What datatype is destination? Arrays are handled differently in Couchbase. Yay!
					if cmd.To[len(cmd.To)-1] == ']' {
						builder = builder.ArrayInsert(cmd.To, x)
					} else {
						builder = builder.Upsert(cmd.To, x, true)
					}
					builder = builder.Remove(cmd.From)
					cmdAdded = true
				}
			case dupblock.ActionSwap:
				var to, from interface{}
				if f, e := b.LookupIn(key.String()).Get(cmd.To).Get(cmd.From).Execute(); e != nil {
					return e
				} else if e = f.Content(cmd.To, &to); e != nil {
					return e
				} else if e = f.Content(cmd.From, &from); e != nil {
					return e
				} else {
					// What datatype is destination? Arrays are handled differently in Couchbase. Yay!
					if cmd.From[len(cmd.From)-1] == ']' {
						builder = builder.ArrayInsert(cmd.From, to)
					} else {
						builder = builder.Upsert(cmd.From, to, true)
					}
					// What datatype is destination? Arrays are handled differently in Couchbase. Yay!
					if cmd.To[len(cmd.To)-1] == ']' {
						builder = builder.ArrayInsert(cmd.To, from)
					} else {
						builder = builder.Upsert(cmd.To, from, true)
					}
					cmdAdded = true
				}
			case dupblock.ActionSetKey:
				if cmdAdded {
					if _, e := builder.Execute(); e != nil {
						return e
					}
					cmdAdded = false
				}
				if key, e = document.NewKeyFromString(cmd.To); e != nil {
					return e
				} else {
					builder = b.MutateIn(key.String(), 0, 0)
				}
			default:
				return dupblock.ErrUnhandledAction
			}
		}

		if cmdAdded {
			if _, e := builder.Execute(); e != nil {
				return convertCouchbaseError(e, key.String())
			}
		}
	}

	return nil
}

func (d *DocumentStore) data() (result *gocb.Bucket, err error) {
	d.mutex.Lock()
	result = d.dataBucket

	if result == nil {
		if result, err = d.owner.cluster.OpenBucket(d.name, ""); err == nil {
			result.SetTranscoder(&DefaultTranscoder)
			d.dataBucket = result
		}
	} else {
		result = d.dataBucket
	}
	d.mutex.Unlock()
	return
}

func (d *DocumentStore) stream() (result *gocb.StreamingBucket, err error) {
	d.mutex.Lock()
	result = d.streamBucket

	if result == nil {
		if result, err = d.owner.cluster.OpenStreamingBucket(
			uuid.New().String(),
			d.name,
			"",
			gocbcore.DcpOpenFlagNoValue,
		); err == nil {
			d.streamBucket = result
		}
	} else {
		result = d.streamBucket
	}
	d.mutex.Unlock()
	return
}

func (d *DocumentStore) watcher() (result *BucketWatcher, err error) {
	d.mutex.Lock()

	result = d.bucketWatcher
	if result == nil {
		// Want to call d.stream(), but do not want to try and lock mutex again.
		stream := d.streamBucket
		if stream == nil {
			if stream, err = d.owner.cluster.OpenStreamingBucket(
				uuid.New().String(),
				d.name,
				"",
				gocbcore.DcpOpenFlagNoValue,
			); err == nil {
				d.streamBucket = stream
			}
		}

		if err == nil {
			result = NewBucketWatcher(stream, d.logger)
			d.bucketWatcher = result
		}
	}

	d.mutex.Unlock()
	return
}

func (d *DocumentStore) AddWatcher(key document.Key, watcher document.DocumentWatcher) {
	if b, err := d.watcher(); err == nil {
		b.WatchDocument(key, watcher)
	}
}

func (d *DocumentStore) RemoveWatcher(key document.Key, watcher document.DocumentWatcher) {
	if b, err := d.watcher(); err == nil {
		b.StopWatchingDocument(key, watcher)
	}
}

func (d *DocumentStore) WatchersForKey(key document.Key) (watchers []document.DocumentWatcher) {
	if b, err := d.watcher(); err == nil {
		watchers = b.WatchersForKey(key)
	}
	return
}

func (d *DocumentStore) PushBack(key document.Key, path string, item interface{}) error {
	if b, e := d.data(); e != nil {
		return e
	} else {
		builder := b.MutateIn(key.String(), 0, 0)
		builder = builder.ArrayAppend(
			path,
			item,
			false,
		)
		if _, e := builder.Execute(); e != nil {
			return convertCouchbaseError(e, key.String())
		} else {
			return nil
		}
	}
}

// Touch touches a document, specifying a new expiry time for it.
func (d *DocumentStore) Touch(key document.Key, opts ...document.Option) (document.Version, error) {
	if o, e := document.NewOptions(opts...); e != nil {
		return document.NoVersion, e
	} else if b, e := d.data(); e != nil {
		return document.NoVersion, convertCouchbaseError(e, "")
	} else {
		cas := gocb.Cas(o.Version)
		if o.AnyVersion {
			cas = 0
		}
		if version, e := b.Touch(key.String(), cas, cbExpiry(o.TTL)); e != nil {
			return document.NoVersion, convertCouchbaseError(e, key.String())
		} else {
			return document.Version(version), nil
		}
	}
}

// Incr increments a numeric field by the provided amount.
// Incr(key Key, path string, amount int32) error
func (d *DocumentStore) Incr(key document.Key, path string, amount int32) error {
	return nil
}
