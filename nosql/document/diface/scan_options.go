package diface

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"moke-kit/nosql/document/key"
	errors2 "moke-kit/nosql/nerrors"
)

const (
	ScanLimitNone    = 0
	ScanLimitDefault = 1

	ScanTimeoutNone    = 0
	ScanTimeoutDefault = 10 * time.Second
)

type ScanQueryType int

const (
	// high level, what kind of a scan is this query?
	ScanTypeUnset ScanQueryType = iota
	ScanTypeKeyValue
	ScanTypeKeyLike
	ScanTypeNOOP
	ScanTypeRegex
	ScanTypeNum
)

type ScanQueryOp int

const (
	// lower level, what operator is being used in the query?
	ScanOpEquals ScanQueryOp = iota
	ScanOpLessThan
	ScanOpGreaterThan
)

type ScanOptions struct {
	Query []ScanQuery

	Offset  int
	Limit   int
	Timeout time.Duration
}

type KeyValue struct {
	Index string
	Value interface{}
}

type ScanQuery struct {
	ScanType      ScanQueryType
	ScanOperation ScanQueryOp

	KeyValue KeyValue
	Regex    string
}

type ScanOption func(o *ScanOptions) error

// NB: consider making a way to request that defaults to unlimited?
func NewScanOptions(opts ...ScanOption) (options ScanOptions, err error) {
	o := &options
	o.Limit = ScanLimitDefault
	o.Timeout = ScanTimeoutDefault
	for _, opt := range opts {
		if err = opt(o); err != nil {
			break
		}
	}
	return
}

func WithLimit(limit int) ScanOption {
	return func(o *ScanOptions) error {
		if limit < 0 {
			return errors.Wrap(errors2.ErrInternal, "negative limit specified")
		} else {
			o.Limit = limit
		}
		return nil
	}
}

func WithNoLimit() ScanOption {
	return WithLimit(ScanLimitNone)
}

// For paginating query results, use in combination with an appropriate limit.
// NB: Offsets are not reliable in mock implementation because mock ordering is nondeterministic
func WithOffset(offset int) ScanOption {
	return func(o *ScanOptions) error {
		if offset < 0 {
			return errors.Wrap(errors2.ErrInternal, "negative offset specified")
		} else {
			o.Offset = offset
		}
		return nil
	}
}

// MatchAny does not perform any filter on the query
func MatchAny() ScanOption {
	return func(o *ScanOptions) error {
		o.Query = append(o.Query, ScanQuery{
			ScanType: ScanTypeNOOP,
		})
		return nil
	}
}

// MatchKeyValue packs up a keys/value search pair when scanning for a document.
func MatchKeyValue(idx string, val string) ScanOption {
	return func(o *ScanOptions) error {
		if idx == "" {
			return errors.Wrap(key.ErrInvalidKeyFormat, "empty keys provided")
		} else {
			o.Query = append(o.Query, ScanQuery{
				ScanType: ScanTypeKeyValue,
				KeyValue: KeyValue{Index: idx, Value: val},
			})
		}
		return nil
	}
}

// MatchKeyLike packs up a keys/value search pair when scanning for a document - for a sql style "like" match.
func MatchKeyLike(idx string, val string) ScanOption {
	return func(o *ScanOptions) error {
		if idx == "" {
			return errors.Wrap(key.ErrInvalidKeyFormat, "empty keys provided")
		} else if len(val) < 2 {
			return errors.Wrap(errors2.ErrInvalidScanValue, "wildcard query string is too short")
		} else {
			// Keys permit the use of periods so they need to be escaped
			// before building the regular expression such that they are
			// interpreted literally.
			val = strings.ReplaceAll(val, ".", `\.`)

			var pre, suf string
			var idxstart, idxend int
			if val[:1] == "%" {
				idxstart = 1
				pre = "^.*?"
			} else {
				idxstart = 0
				pre = "^"
			}
			if val[len(val)-1:] == "%" {
				idxend = len(val) - 1
				suf = ".*?$"
			} else {
				idxend = len(val)
				suf = "$"
			}

			o.Query = append(o.Query, ScanQuery{
				ScanType: ScanTypeKeyLike,
				KeyValue: KeyValue{Index: idx, Value: val},
				Regex:    fmt.Sprint(pre, val[idxstart:idxend], suf),
			})
		}
		return nil
	}
}

func MatchRegex(regex string) ScanOption {
	return func(o *ScanOptions) error {
		if regex == "" {
			return errors.Wrap(errors2.ErrInternal, "empty regex provided")
		} else {
			o.Query = append(o.Query, ScanQuery{
				ScanType: ScanTypeRegex,
				Regex:    regex,
			})
		}
		return nil
	}
}

func MatchNumber(idx string, op ScanQueryOp, val float64) ScanOption {
	return func(o *ScanOptions) error {
		if idx == "" {
			return errors.Wrap(key.ErrInvalidKeyFormat, "empty keys provided")
		} else {
			o.Query = append(o.Query, ScanQuery{
				ScanType:      ScanTypeNum,
				ScanOperation: op,
				KeyValue:      KeyValue{Index: idx, Value: val},
			})
		}
		return nil
	}
}

// Use WithTimeout to set the expiration time of a supported operation.
func WithTimeout(timeout time.Duration) ScanOption {
	return func(o *ScanOptions) error {
		o.Timeout = timeout
		return nil
	}
}

func WithNoTimeout() ScanOption {
	return WithTimeout(ScanTimeoutNone)
}
