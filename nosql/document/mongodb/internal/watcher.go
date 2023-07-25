package internal

//type BucketWatcher struct {
//	sync.Mutex
//
//	logger            *zap.Logger
//	coll              *mongo.Collection
//	vBucketWatcherMap map[document.Key]*VBucketWatcher
//}
//
//func NewBucketWatcher(collection *mongo.Collection, logger *zap.Logger) *BucketWatcher {
//	result := &BucketWatcher{
//		vBucketWatcherMap: map[document.Key]*VBucketWatcher{},
//		coll:              collection,
//		logger:            logger,
//	}
//
//	return result
//}
//
//func (w *BucketWatcher) WatchDocument(key document.Key, handler document.DocumentWatcher) {
//	w.logger.Info(
//		"WatchDocument",
//		zap.String("keys", key.String()),
//		zap.String("keys-bytes", string(key.Bytes())),
//	)
//
//	w.Lock()
//	defer w.Unlock()
//	{
//		var watcher *VBucketWatcher
//		var ok bool
//		if watcher, ok = w.vBucketWatcherMap[key]; !ok {
//			watcher = NewVBucketWatcher(key, w)
//			w.vBucketWatcherMap[key] = watcher
//		}
//		watcher.AddWatcher(key, handler)
//	}
//}
//
//func (w *BucketWatcher) StopWatchingDocument(key document.Key, handler document.DocumentWatcher) {
//	w.Lock()
//	{
//		if watcher, ok := w.vBucketWatcherMap[key]; ok {
//			if watcher.RemoveWatcher(key, handler) {
//				watcher.Close()
//			}
//
//			w.logger.Info(
//				"StopWatchingDocument",
//				zap.String("keys", key.String()),
//				zap.String("keys-bytes", string(key.Bytes())),
//			)
//		} else {
//			w.logger.Info(
//				"StopWatchingDocument - no watcher found",
//				zap.String("keys", key.String()),
//				zap.String("keys-bytes", string(key.Bytes())),
//			)
//		}
//	}
//	w.Unlock()
//}
//
//func (w *BucketWatcher) WatchersForKey(key document.Key) (watchers []document.DocumentWatcher) {
//	w.Lock()
//	{
//		if watcher, ok := w.vBucketWatcherMap[key]; ok {
//			watchers = watcher.WatchersForKey(key)
//		}
//	}
//	w.Unlock()
//	return
//}
