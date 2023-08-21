package nats

import (
	"github.com/gstones/moke-kit/mq/logic"
)

type message struct {
	id    string
	topic string
	data  []byte
	vPtr  interface{}
}

func NewMessage(id string, topic string, data []byte, vPtr interface{}) logic.Message {
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

func (m *message) VPtr() (vPtr interface{}) {
	return m.vPtr
}
