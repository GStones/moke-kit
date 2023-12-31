package noptions

import (
	"reflect"
	"time"

	"github.com/gstones/moke-kit/orm/nerrors"
)

type Version = int64

const (
	NoVersion Version = 0
)

// Options contains all the various options that the provided WithXyz functions construct.
type Options struct {
	Version         Version
	AnyVersion      bool
	TTL             time.Duration
	Source          any
	Sources         map[string]any
	Destination     any
	DestinationList []any
}

// Option is a closure that updates Options.
type Option func(o *Options) error

// NewOptions constructs an Options struct from the provided Option closures and returns it.
func NewOptions(opts ...Option) (options Options, err error) {
	o := &options
	o.Version = NoVersion
	for _, opt := range opts {
		if err = opt(o); err != nil {
			break
		}
	}
	return
}

// WithVersion provides a Version value option.  This is optionally used when updating documents to ensure the user has
// the correct version before updating.
func WithVersion(v Version) Option {
	return func(o *Options) error {
		if o.AnyVersion {
			return nerrors.ErrAnyVersionConflict
		} else {
			o.Version = v
			return nil
		}
	}
}

// Use WithAnyVersion to update documents without caring what their version is.
func WithAnyVersion() Option {
	return func(o *Options) error {
		if o.Version != NoVersion {
			return nerrors.ErrAnyVersionConflict
		} else {
			o.AnyVersion = true
			return nil
		}
	}
}

// Use WithTTL to set the expiration time of a nosql during an operation.
func WithTTL(ttl time.Duration) Option {
	return func(o *Options) error {
		o.TTL = ttl
		return nil
	}
}

// WithSource provides an interface to source data when updating a nosql.
func WithSource(src any) Option {
	return func(o *Options) (err error) {
		o.Source = src
		return
	}
}

func WithMultipleSource(src map[string]any) Option {
	return func(o *Options) (err error) {
		o.Sources = src
		return
	}
}

// WithDestination provides an interface for receiving data when getting a nosql.
func WithDestination(dst any) Option {
	return func(o *Options) error {
		if dst == nil {
			return nerrors.ErrDestIsNil
		} else if reflect.TypeOf(dst).Kind() != reflect.Ptr {
			return nerrors.ErrDestMustBePointer
		} else {
			o.Destination = dst
			return nil
		}
	}
}

// WithDestinationList provides a interface list for receiving data when getting a nosql.
func WithDestinationList(dst []any) Option {
	return func(o *Options) error {
		if dst == nil {
			return nerrors.ErrDestIsNil
		} else {
			o.DestinationList = dst
			return nil
		}
	}
}
