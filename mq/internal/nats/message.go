package nats

import (
	"github.com/gstones/moke-kit/mq/miface"
)

type message struct {
	id    string
	topic string
	data  []byte
	vPtr  any
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
