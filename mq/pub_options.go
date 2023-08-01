package mq

import (
	"encoding/json"
	"time"
)

// PubOptions contains all the various options that the provided WithXyz functions construct.
type PubOptions struct {
	Data  []byte
	Delay time.Duration
}

// PubOption is a closure that updates PubOptions.
type PubOption func(o *PubOptions) error

// NewPubOptions constructs an PubOptions struct from the provided PubOption closures and returns it.
func NewPubOptions(opts ...PubOption) (options PubOptions, err error) {
	o := &options
	for _, opt := range opts {
		if err = opt(o); err != nil {
			break
		}
	}
	return
}

// Use WithBytes to set mq's message []byte data payload directly
func WithBytes(data []byte) PubOption {
	return func(o *PubOptions) error {
		if len(o.Data) != 0 {
			return ErrDataAlreadySet
		} else {
			o.Data = data
			return nil
		}
	}
}

// Use WithJson to set the PubOptions' Data field with a JSON object
func WithJson(data interface{}) PubOption {
	return func(o *PubOptions) (err error) {
		if len(o.Data) != 0 {
			return ErrDataAlreadySet
		} else {
			o.Data, err = json.Marshal(data)
			return
		}
	}
}

// Use WithEncoder to set the PubOptions' Data field with with encoded data
func WithEncoder(encoder Encoder, topic string, data interface{}) PubOption {
	return func(o *PubOptions) error {
		if len(o.Data) != 0 {
			return ErrDataAlreadySet
		} else if data, err := encoder.Encode(NamespaceTopic(topic), data); err != nil {
			return err
		} else {
			o.Data = data
			return nil
		}
	}
}

// Use WithDelay to set the value to defer the message. Deferring is to be done by the message queue, not mq.
func WithDelay(deferAmt time.Duration) PubOption {
	return func(o *PubOptions) error {
		o.Delay = deferAmt
		return nil
	}
}
