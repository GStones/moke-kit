package logic

type MessageQueue interface {
	Subscribe(topic string, handler SubResponseHandler, opts ...SubOption) (Subscription, error)
	Publish(topic string, opts ...PubOption) error
}

type Subscription interface {
	IsValid() bool
	Unsubscribe() error
}

type Message interface {
	ID() string
	Topic() string
	Data() []byte
	VPtr() (vPtr any)
}

type Encoder interface {
	Encode(subject string, v any) ([]byte, error)
}

type Decoder interface {
	Decode(subject string, data []byte, vPtr any) error
}

type Codec interface {
	Encoder
	Decoder
}

// ValuePtrFactory is used during optional subscription decoding.
// A ValuePtrFactory produces the value pointer populated by Decode()
type ValuePtrFactory interface {
	NewVPtr() (vPtr any)
}
