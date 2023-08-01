package internal

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"sync"

	"moke-kit/mq"
)

type Consumer struct {
	sync.RWMutex
	c           *nsq.Consumer
	decoder     mq.Decoder
	vPtrFactory mq.ValuePtrFactory
	handler     mq.SubResponseHandler
	valid       bool
	topic       string
}

func NewNsqConsumer(topic string,
	address string,
	handler mq.SubResponseHandler,
	decoder mq.Decoder,
	vPtrFactory mq.ValuePtrFactory,
) (*Consumer, error) {
	c := &Consumer{
		decoder:     decoder,
		handler:     handler,
		topic:       topic,
		vPtrFactory: vPtrFactory,
	}
	channel := "nsq"

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}

	consumer.AddHandler(c)
	err = consumer.ConnectToNSQLookupd(address)
	if err != nil {
		return nil, err
	}
	c.c = consumer
	c.valid = true

	return c, nil
}

func (c *Consumer) HandleMessage(message *nsq.Message) error {
	id := fmt.Sprintf("%x", message.ID)
	if c.decoder != nil && c.vPtrFactory != nil {
		vPtrMessage := c.vPtrFactory.NewVPtr()

		if err := c.decoder.Decode(c.topic, message.Body, vPtrMessage); err != nil {
			c.handler(nil, err)
		} else {
			mqMsg := NewMessage(id, c.topic, nil, vPtrMessage)
			c.handler(mqMsg, nil)
		}
	} else {
		mqMsg := NewMessage(id, c.topic, message.Body, nil)

		c.handler(mqMsg, nil)
	}

	return nil
}

func (c *Consumer) IsValid() bool {
	c.RLock()
	defer c.RUnlock()

	return c.valid
}

func (c *Consumer) Unsubscribe() error {
	c.Lock()
	if c.valid {
		c.valid = false
		c.c.Stop()
		<-c.c.StopChan
	}
	c.Unlock()

	return nil
}
