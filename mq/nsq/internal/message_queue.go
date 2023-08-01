package internal

import (
	"go.uber.org/zap"
)

type messageQueue struct {
	logger *zap.Logger

	consumerAddress string
	producerAddress string

	producer *nsq.Producer
}

func NewMessageQueue(logger *zap.Logger, consumerAddress string, producerAddress string) (mq *messageQueue, err error) {
	config := nsq.NewConfig()

	if e := validateEndpoint(producerAddress, config); e != nil {
		err = e
	} else if p, e := nsq.NewProducer(producerAddress, config); e != nil {
		err = e
	} else {
		mq = &messageQueue{
			logger:          logger,
			consumerAddress: consumerAddress,
			producerAddress: producerAddress,
			producer:        p,
		}
	}

	return
}

func (m *messageQueue) Subscribe(
	topic string,
	handler mq.SubResponseHandler,
	opts ...mq.SubOption,
) (mq.Subscription, error) {
	if o, err := mq.NewSubOptions(opts...); err != nil {
		return nil, err
	} else if o.DeliverySemantics == mq.AtMostOnce {
		return nil, mq.ErrAtMostOnceUnsupported
	} else {
		return NewNsqConsumer(topic, m.consumerAddress, handler, o.Decoder, o.VPtrFactory)
	}
}

func (m *messageQueue) Publish(topic string, opts ...mq.PubOption) error {
	if options, err := mq.NewPubOptions(opts...); err != nil {
		m.logger.Error(
			"Publish error:",
			zap.Error(err),
		)
		return err
	} else if options.Delay != 0 {
		return m.producer.DeferredPublish(topic, options.Delay, options.Data)
	} else {
		return m.producer.Publish(topic, options.Data)
	}
}

func validateEndpoint(address string, config *nsq.Config) error {
	conn := nsq.NewConn(address, config, &connDelegate{})
	if _, err := conn.Connect(); err != nil {
		return err
	} else {
		conn.Close()
	}

	return nil
}

type connDelegate struct{}

func (cd *connDelegate) OnResponse(c *nsq.Conn, data []byte) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnError(c *nsq.Conn, data []byte) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnMessage(c *nsq.Conn, msg *nsq.Message) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnMessageFinished(c *nsq.Conn, msg *nsq.Message) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnMessageRequeued(c *nsq.Conn, msg *nsq.Message) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnBackoff(*nsq.Conn) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnContinue(*nsq.Conn) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnResume(*nsq.Conn) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnIOError(c *nsq.Conn, err error) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnHeartbeat(*nsq.Conn) {
	// Do nothing. This is here just to test URL.
}

func (cd *connDelegate) OnClose(*nsq.Conn) {
	// Do nothing. This is here just to test URL.
}
