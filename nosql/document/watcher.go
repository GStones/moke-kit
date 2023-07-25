package document

import (
	"sync"

	"go.uber.org/zap"
)

// DocumentWatcherManager manages document watchers.
type DocumentWatcherManager struct {
	mutex         sync.RWMutex
	keyToWatchers map[Key][]DocumentWatcher
	Logger        *zap.Logger
}

// AddWatcher adds a DocumentWatcher to the given keys.  Whenever changes to the document occur the watcher is told.
func (d *DocumentWatcherManager) AddWatcher(key Key, watcher DocumentWatcher) {
	d.Logger.Info("AddWatcher", zap.String("Key", key.String()))
	d.mutex.Lock()
	{
		if d.keyToWatchers == nil {
			d.keyToWatchers = map[Key][]DocumentWatcher{
				key: {watcher},
			}
		} else {
			d.keyToWatchers[key] = append(d.keyToWatchers[key], watcher)
		}
	}
	d.mutex.Unlock()
}

// RemoveWatcher removes a DocumentWatcher from the given keys.  It's the cleanup pair for AddWatcher.
func (d *DocumentWatcherManager) RemoveWatcher(key Key, watcher DocumentWatcher) {
	d.Logger.Info("RemoveWatcher", zap.String("Key", key.String()))
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.keyToWatchers != nil {
		if watchers, ok := d.keyToWatchers[key]; ok {
			for i, ew := range watchers {
				if ew == watcher {
					d.keyToWatchers[key] = append(watchers[:i], watchers[i+1:]...)
					if len(d.keyToWatchers[key]) == 0 {
						delete(d.keyToWatchers, key)
					}
				}
			}
		}
	}
}

// OnDocumentChanged dispatches a change notification for all watchers of keys.
func (d *DocumentWatcherManager) OnDocumentChanged(key Key, changes interface{}) {
	d.forEachWatcher(key, func(dw DocumentWatcher) {
		dw.OnDocumentChanged(key, changes)
	})
}

// OnAllDocumentsChanged dispatches a change notification for all watchers.
func (d *DocumentWatcherManager) OnAllDocumentsChanged() {
	d.forEachKey(d.OnDocumentChanged)
}

// OnDocumentExpired dispatches an expired notification for all watchers of keys.
func (d *DocumentWatcherManager) OnDocumentExpired(key Key) {
	d.forEachWatcher(key, func(dw DocumentWatcher) {
		dw.OnDocumentExpired(key)
	})
}

// OnDocumentDeleted dispatches a deleted notification for all watchers of keys.
func (d *DocumentWatcherManager) OnDocumentDeleted(key Key) {
	d.forEachWatcher(key, func(dw DocumentWatcher) {
		dw.OnDocumentDeleted(key)
	})
}

// WatchersForKey returns a new slice of all DocumentWatcher instances for a keys.
func (d *DocumentWatcherManager) WatchersForKey(key Key) (watchers []DocumentWatcher) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if d.keyToWatchers != nil {
		for _, w := range d.keyToWatchers[key] {
			watchers = append(watchers, w)
		}
	}
	return
}

// Keys returns a new slice of all keys being watched by this manager.
func (d *DocumentWatcherManager) Keys() (keys []Key) {
	d.mutex.RLock()
	{
		for k := range d.keyToWatchers {
			keys = append(keys, k)
		}
	}
	d.mutex.RUnlock()
	return
}

// forEachWatcher iterates over all watchers for a keys and invokes a callback.
func (d *DocumentWatcherManager) forEachWatcher(key Key, cb func(w DocumentWatcher)) (ok bool) {
	for _, watcher := range d.WatchersForKey(key) {
		cb(watcher)
		ok = true
	}
	return
}

// forEachWatcher iterates over all keys and invokes a callback.
func (d *DocumentWatcherManager) forEachKey(cb func(k Key, v interface{})) {
	keys := d.Keys()
	for _, k := range keys {
		cb(k, nil)
	}
}
