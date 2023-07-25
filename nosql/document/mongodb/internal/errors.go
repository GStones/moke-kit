package internal

import (
	errors2 "github.com/gstones/platform/services/common/nosql/errors"
	"github.com/pkg/errors"
)

// Couchbase specific errors.
var (
	ErrCompressionNotSupported  = errors.New("ErrCompressionNotSupported")
	ErrBinaryTypeNotSupported   = errors.New("ErrBinaryTypeNotSupported")
	ErrStringTypeNotSupported   = errors.New("ErrStringTypeNotSupported")
	ErrUnknownValueType         = errors.New("ErrUnknownValueType")
	ErrInvalidDestination       = errors.New("ErrInvalidDestination")
	ErrN1QLResultIncomplete     = errors.New("ErrN1QLResultIncomplete")
	ErrN1QLResultFormatMismatch = errors.New("ErrN1QLResultFormatMismatch")
	ErrDocumentStoreIsNil       = errors.New("ErrDocumentStoreIsNil")
)

func errInternal(err error) error {
	return errors.Wrap(errors2.ErrInternal, err.Error())
}
