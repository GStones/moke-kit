package internal

import (
	"sync"

	"github.com/gstones/platform/services/common/mq"
)

// This struct exists to support multiple at-most-once subscriptions
// who share a topic but may or may not differ in their respective groupId.
// All readerGroups under an idGroup share a single topic in common.
type idGroup struct {
	mutex               sync.Mutex
	activeAtMostOnceRGs map[string]*readerGroup
}

func NewIdGroup() *idGroup {
	return &idGroup{activeAtMostOnceRGs: make(map[string]*readerGroup)}
}

func (i *idGroup) OpenReaderGroup(groupId string) (*readerGroup, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if rg, ok := i.activeAtMostOnceRGs[groupId]; ok {
		return rg, nil
	} else {
		return nil, mq.ErrReaderGroupNotFound
	}
}

func (i *idGroup) AddReaderGroup(groupId string, rg *readerGroup) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if _, ok := i.activeAtMostOnceRGs[groupId]; ok {
		return mq.ErrReaderGroupAlreadyExists
	} else {
		i.activeAtMostOnceRGs[groupId] = rg
		return nil
	}
}
