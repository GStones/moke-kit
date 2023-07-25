package internal

import (
	"github.com/gstones/platform/services/common/nosql/document"
	"sync"

	"github.com/snichols/gocb"
	"go.uber.org/zap"
)

// SNICHOLS: https://github.com/couchbaselabs/dcp-documentation

type BucketWatcher struct {
	sync.Mutex

	logger            *zap.Logger
	stream            *gocb.StreamingBucket
	vBucketWatcherMap map[uint16]*vBucketWatcher
}

func NewBucketWatcher(streamBucket *gocb.StreamingBucket, logger *zap.Logger) *BucketWatcher {
	result := &BucketWatcher{
		vBucketWatcherMap: map[uint16]*vBucketWatcher{},
		stream:            streamBucket,
		logger:            logger,
	}

	return result
}

func (w *BucketWatcher) WatchDocument(key document.Key, handler document.DocumentWatcher) {
	vbId := w.stream.IoRouter().KeyToVbucket(key.Bytes())

	w.logger.Info(
		"WatchDocument",
		zap.String("keys", key.String()),
		zap.String("keys-bytes", string(key.Bytes())),
		zap.Uint16("vbId", vbId),
	)

	w.Lock()
	{
		var watcher *vBucketWatcher
		var ok bool
		if watcher, ok = w.vBucketWatcherMap[vbId]; !ok {
			watcher = NewVBucketWatcher(vbId, w)
			w.vBucketWatcherMap[vbId] = watcher
		}
		watcher.AddWatcher(key, handler)
	}
	w.Unlock()

}

func (w *BucketWatcher) StopWatchingDocument(key document.Key, handler document.DocumentWatcher) {
	vbId := w.stream.IoRouter().KeyToVbucket(key.Bytes())

	w.Lock()
	{
		if watcher, ok := w.vBucketWatcherMap[vbId]; ok {
			watcher.RemoveWatcher(key, handler)
			w.logger.Info(
				"StopWatchingDocument",
				zap.String("keys", key.String()),
				zap.String("keys-bytes", string(key.Bytes())),
				zap.Uint16("vbId", vbId),
			)
		} else {
			w.logger.Info(
				"StopWatchingDocument - no watcher found",
				zap.String("keys", key.String()),
				zap.String("keys-bytes", string(key.Bytes())),
				zap.Uint16("vbId", vbId),
			)
		}
	}
	w.Unlock()
}

func (w *BucketWatcher) WatchersForKey(key document.Key) (watchers []document.DocumentWatcher) {
	vbId := w.stream.IoRouter().KeyToVbucket(key.Bytes())

	w.Lock()
	{
		if watcher, ok := w.vBucketWatcherMap[vbId]; ok {
			watchers = watcher.WatchersForKey(key)
		}
	}
	w.Unlock()
	return
}
