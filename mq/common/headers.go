package common

type Header string

const (
	KafkaHeader Header = "kafka://"
	NatsHeader  Header = "nats://"
	NsqHeader   Header = "nsq://"
	LocalHeader Header = "local://"
)

func (h Header) CreateTopic(topic string) string {
	return string(h) + topic
}
