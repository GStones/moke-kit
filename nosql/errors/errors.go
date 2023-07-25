package errors

import (
	"errors"
)

var (
	ErrAnyVersionConflict  = errors.New("ErrAnyVersionConflict")
	ErrDestIsNil           = errors.New("ErrDestIsNil")
	ErrDestMustBePointer   = errors.New("ErrDestMustBePointer")
	ErrKeyExists           = errors.New("ErrKeyExists")
	ErrKeyNotFound         = errors.New("ErrKeyNotFound")
	ErrNoResults           = errors.New("ErrNoResults")
	ErrNoSuchDocumentStore = errors.New("ErrNoSuchDocumentStore")
	ErrDocumentStoreIsNil  = errors.New("ErrDocumentStoreIsNil")
	ErrSourceIsNil         = errors.New("ErrSourceIsNil")
	ErrTooManyRetries      = errors.New("ErrTooManyRetries")
	ErrUnknownDUPAction    = errors.New("ErrUnknownDUPAction")
	ErrUpdateLogicFailed   = errors.New("ErrUpdateLogicFailed")
	ErrVersionMismatch     = errors.New("ErrVersionMismatch")
	ErrInvalidVersioning   = errors.New("ErrInvalidVersioning")
	ErrInvalidStoreName    = errors.New("ErrInvalidStoreName")
	ErrNotImplemented      = errors.New("ErrNotImplemented")
	ErrDriverFailure       = errors.New("ErrDriverFailure")
	ErrMalformedData       = errors.New("ErrMalformedData")
	ErrUnexpectedScanType  = errors.New("ErrUnexpectedScanType")
	ErrNoScanType          = errors.New("ErrNoScaSnType")
	ErrInvalidScanValue    = errors.New("ErrInvalidScanValue")
	ErrInternal            = errors.New("ErrInternal")
	ErrNotFloat64          = errors.New("ErrNotFloat64")
)

func IsNotFound(err error) bool {
	return err == ErrKeyNotFound
}
