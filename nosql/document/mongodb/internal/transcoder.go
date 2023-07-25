package internal

import (
	"github.com/gstones/platform/services/common/jsonx"

	"gopkg.in/couchbase/gocbcore.v7"
)

var (
	DefaultTranscoder     = Transcoder{}
	UncompressedJsonFlags = gocbcore.EncodeCommonFlags(gocbcore.JsonType, gocbcore.NoCompression)
)

// Transcoder overrides Couchbase's default transcoder for our document store use cases.  No binary / string support
// and using the faster jsoniter library.
type Transcoder struct {
}

// Decodes retrieved bytes into a Go type.
func (t Transcoder) Decode(data []byte, flags uint32, dest interface{}) error {
	valueType, compression := gocbcore.DecodeCommonFlags(flags)

	// @TODO: SNICHOLS: support compression one day?
	if compression != gocbcore.NoCompression {
		return errInternal(ErrCompressionNotSupported)
	}

	switch valueType {
	case gocbcore.BinaryType:
		return errInternal(ErrBinaryTypeNotSupported)

	case gocbcore.StringType:
		return errInternal(ErrStringTypeNotSupported)

	case gocbcore.JsonType:
		return jsonx.Unmarshal(data, dest)

	default:
		return errInternal(ErrUnknownValueType)
	}
}

// Encodes a Go type into bytes for storage.
func (t Transcoder) Encode(value interface{}) (data []byte, flags uint32, err error) {
	flags = UncompressedJsonFlags
	data, err = jsonx.Marshal(value)
	return
}
