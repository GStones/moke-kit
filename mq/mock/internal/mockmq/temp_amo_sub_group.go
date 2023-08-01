package mockmq

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gstones/platform/services/common/mq"
)

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	randMutex = sync.Mutex{}
)

// A temporary SubGroup created on per-publish basis to keep track of
// subscriptions eligible for at-most-once message delivery
type tempAtMostOnceSubGroup struct {
	sync.Mutex
	queuedSubs map[string][]Subscription
}

func NewTempAtMostOnceSubGroup() TempAtMostOnceSubGroup {
	return &tempAtMostOnceSubGroup{queuedSubs: make(map[string][]Subscription)}
}

func (t *tempAtMostOnceSubGroup) AddSubscriptions(subs ...Subscription) {
	t.Lock()

	for _, sub := range subs {
		if queueName := sub.Queue(); queueName == "" {
			continue
		} else if _, ok := t.queuedSubs[queueName]; ok {
			t.queuedSubs[queueName] = append(t.queuedSubs[queueName], sub)
		} else {
			t.queuedSubs[queueName] = make([]Subscription, 0)
			t.queuedSubs[queueName] = append(t.queuedSubs[queueName], sub)
		}
	}
	t.Unlock()
}

// For each queue group, only one Subscriber will be sent the published message
func (t *tempAtMostOnceSubGroup) SendMessage(topic string, data []byte) {

	t.Lock()

	message := Message{Topic: topic, Data: data}

	for _, namedQueue := range t.queuedSubs {
	queueIterate:
		for {
			if len(namedQueue) == 0 {
				break queueIterate
			} else {
				randMutex.Lock()
				randomIndex := random.Intn(len(namedQueue))
				randMutex.Unlock()

				if err := namedQueue[randomIndex].SendMessage(message); err != mq.ErrInvalidSubscription {
					break queueIterate
				} else {
					namedQueue = append(namedQueue[:randomIndex], namedQueue[randomIndex+1:]...)
				}
			}
		}
	}
	t.Unlock()
}

func (t *tempAtMostOnceSubGroup) Subscriptions() []Subscription {
	t.Lock()

	var subs []Subscription
	for _, namedQueue := range t.queuedSubs {
		subs = append(subs, namedQueue...)
	}

	t.Unlock()

	return subs
}
