package message

import (
	wmsg "github.com/ThreeDotsLabs/watermill/message"

	"github.com/gstones/moke-kit/mq/miface"
)

type message struct {
	id    string
	topic string
	data  []byte
	vPtr  any
}

// Msg2Message converts a watermill message to a message
func Msg2Message(topic string, msg *wmsg.Message) miface.Message {
	return NewMessage(msg.UUID, topic, msg.Payload, nil)
}

func NewMessage(id string, topic string, data []byte, vPtr any) miface.Message {
	return &message{
		id:    id,
		topic: topic,
		data:  data,
		vPtr:  vPtr,
	}
}

func (m *message) ID() string {
	return m.id
}

func (m *message) Topic() string {
	return m.topic
}

func (m *message) Data() []byte {
	return m.data
}

func (m *message) VPtr() (vPtr any) {
	return m.vPtr
}
