package common

type Header string

const (
	NatsHeader  Header = "nats://"
	LocalHeader Header = "local://"

	// KafkaHeader and NsqHeader are retained for topic-string compatibility.
	// moke-kit does not ship Kafka or NSQ backends; those prefixes fail as unsupported.
	KafkaHeader Header = "kafka://"
	NsqHeader   Header = "nsq://"
)

func (h Header) CreateTopic(topic string) string {
	return string(h) + topic
}
