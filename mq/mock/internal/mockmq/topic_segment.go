package mockmq

import (
	"strings"
	"sync"
)

// We organize subscriptions into a tree structure made up of topicSegments
// This allows us to recursively walk the tree to send Messages to all appropriate subscriptions
type topicSegment struct {
	mutex sync.Mutex

	subs             PermAtLeastOnceSubGroup // Subscriptions whose topic terminates at this node
	tailWildcardSubs PermAtLeastOnceSubGroup // TailWildcard Subscriptions whose topic terminates at this node

	queuedSubs             PermAtMostOnceSubGroup // Queued Subscriptions whose topic terminates at this node
	queuedTailWildcardSubs PermAtMostOnceSubGroup // Queued TailWildcard Subscriptions whose topic terminates at this node

	childSegments map[string]*topicSegment
}

func NewTopicSegment() *topicSegment {
	return &topicSegment{
		subs:                   NewPermSubGroup(),
		tailWildcardSubs:       NewPermSubGroup(),
		queuedSubs:             NewPermSubGroup(),
		queuedTailWildcardSubs: NewPermSubGroup(),
		childSegments:          make(map[string]*topicSegment),
	}
}

// Recursively walk the tree to subscribe at the appropriate topicSegment node,
// taking into account both Wildcard (*) and TailWildcard (>) subscriptions.
func (t *topicSegment) subscribe(topic string, queue string, handler MessageHandler) Subscription {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// First base case -- An exact match was found!
	if topic == "" {
		if queue == "" {
			return t.subs.NewSubscription(queue, handler)
		} else {
			return t.queuedSubs.NewSubscription(queue, handler)
		}

		// Second base case -- This is a TailWildcard subscription
	} else if topic == TailWildcard {
		if segment, ok := t.childSegments[TailWildcard]; ok {
			return segment.tailWildcardSubscribe(queue, handler)
		} else {
			t.childSegments[TailWildcard] = NewTopicSegment()
			return t.childSegments[TailWildcard].tailWildcardSubscribe(queue, handler)
		}

		// Further recursive topic parsing is required
	} else {
		topicSegments := strings.SplitN(topic, TopicSep, 2)
		leftSide := topicSegments[0]
		rightSide := ""

		if len(topicSegments) == 2 {
			rightSide = topicSegments[1]
		}

		if segment, ok := t.childSegments[leftSide]; ok {
			return segment.subscribe(rightSide, queue, handler)
		} else {
			t.childSegments[leftSide] = NewTopicSegment()
			return t.childSegments[leftSide].subscribe(rightSide, queue, handler)
		}
	}
}

func (t *topicSegment) tailWildcardSubscribe(queue string, handler MessageHandler) Subscription {
	if queue == "" {
		return t.tailWildcardSubs.NewSubscription(queue, handler)
	} else {
		return t.queuedTailWildcardSubs.NewSubscription(queue, handler)
	}
}

// Recursively walk the tree to publish the message to all matching subscriptions,
// taking into account both Wildcard (*) and TailWildcard (>) subscriptions.
func (t *topicSegment) publish(
	topic string,
	fullTopic string,
	data []byte,
	syntheticQueueGroup TempAtMostOnceSubGroup,
) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Base case -- An exact match was found!
	if topic == "" {
		t.subs.SendMessage(fullTopic, data)
		syntheticQueueGroup.AddSubscriptions(t.queuedSubs.Subscriptions()...)

		// Further recursive topic parsing is required
	} else {
		topicSegments := strings.SplitN(topic, TopicSep, 2)
		leftSide := topicSegments[0]
		rightSide := ""

		if len(topicSegments) == 2 {
			rightSide = topicSegments[1]
		}

		if segment, ok := t.childSegments[leftSide]; ok {
			segment.publish(rightSide, fullTopic, data, syntheticQueueGroup)
		}

		// If one exists, we always publish() down the `*` Wildcard branch as well
		// For example: `foo.bar.biz` becomes `*.bar.biz`
		if segment, ok := t.childSegments[Wildcard]; ok {
			segment.publish(rightSide, fullTopic, data, syntheticQueueGroup)
		}

		// If this isn't the base case, we handle any possible `>` TailWildcard subscriptions
		// For example: `foo.bar.biz` becomes `>`
		if segment, ok := t.childSegments[TailWildcard]; ok {
			segment.tailWildcardSubs.SendMessage(fullTopic, data)
			syntheticQueueGroup.AddSubscriptions(segment.queuedTailWildcardSubs.Subscriptions()...)
		}
	}
}
