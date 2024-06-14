package nerrors

import (
	"errors"
)

var (
	ErrNotFound           = errors.New("ErrNotFound")
	ErrVersionNotMatch    = errors.New("ErrVersionNotMatch")
	ErrMissingNosqlURL    = errors.New("ErrMissingNosqlURL")
	ErrInvalidNosqlURL    = errors.New("ErrInvalidNosqlURL")
	ErrAnyVersionConflict = errors.New("ErrAnyVersionConflict")
	ErrDestIsNil          = errors.New("ErrDestIsNil")
	ErrDestMustBePointer  = errors.New("ErrDestMustBePointer")
	ErrKeyNotFound        = errors.New("ErrKeyNotFound")
	ErrDocumentStoreIsNil = errors.New("ErrDocumentStoreIsNil")
	ErrSourceIsNil        = errors.New("ErrSourceIsNil")
	ErrTooManyRetries     = errors.New("ErrTooManyRetries")
	ErrUpdateLogicFailed  = errors.New("ErrUpdateLogicFailed")
)
