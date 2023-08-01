package mockmq

import (
	"regexp"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gstones/platform/services/common/mq"
)

const (
	// If changing a reserved character, update regex.Compile() expression string below
	TopicSep     = "."
	Wildcard     = "*"
	TailWildcard = ">"
)

type messageQueue struct {
	logger      *zap.Logger
	mutex       sync.Mutex
	rootSegment *topicSegment
	deployment  string
}

func NewMessageQueue(logger *zap.Logger, deployment string) *messageQueue {
	return &messageQueue{logger: logger, rootSegment: NewTopicSegment(), deployment: deployment}
}

func (m *messageQueue) Subscribe(
	topic string,
	handler MessageHandler,
) (Subscription, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.validateTopicSemantics(topic, true); err != nil {
		return nil, err
	}

	return m.rootSegment.subscribe(topic, "", handler), nil
}

func (m *messageQueue) QueueSubscribe(
	topic string,
	queue string,
	handler MessageHandler,
) (Subscription, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.validateTopicSemantics(topic, true); err != nil {
		return nil, err
	}

	if queue == "" {
		return nil, mq.ErrEmptyQueueValue
	} else {
		return m.rootSegment.subscribe(topic, queue, handler), nil
	}
}

func (m *messageQueue) Publish(topic string, data []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.validateTopicSemantics(topic, false); err != nil {
		return err
	}

	eligibleQueuedSubs := NewTempAtMostOnceSubGroup()
	m.rootSegment.publish(topic, topic, data, eligibleQueuedSubs)
	eligibleQueuedSubs.SendMessage(topic, data)

	return nil
}

func (m *messageQueue) PublishWithDelay(topic string, data []byte, delay time.Duration) error {
	time.AfterFunc(delay, func() {
		if e := m.Publish(topic, data); e != nil {
			m.logger.Error("Error publishing after delay.",
				zap.Error(e),
			)
		}
	})
	return nil
}

func (m *messageQueue) validateTopicSemantics(topic string, allowWildcards bool) (err error) {
	if topic == "" || topic == m.deployment {
		return mq.ErrEmptyTopic
	}

	var r *regexp.Regexp

	if allowWildcards {
		if r, err = regexp.Compile("[>]|(([*]|[^*>.\\s]+)((?:[.](?:[*]|[^*>.\\s]+))*)([.][>]){0,1})"); err != nil {
			return err
		}
	} else {
		if r, err = regexp.Compile("([^*>.\\s]+)([.]([^*>.\\s]+))*"); err != nil {
			return err
		}
	}

	if loc := r.FindStringIndex(topic); len(loc) != 2 {
		return mq.ErrUnexpectedRegexResults
	} else if loc[0] != 0 {
		return mq.ErrTopicValidationFail
	} else if loc[1] != len(topic) {
		return mq.ErrTopicValidationFail
	} else {
		return nil
	}
}
