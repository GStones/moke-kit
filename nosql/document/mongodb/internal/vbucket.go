package internal

import (
	"context"
	"fmt"
	"github.com/gstones/platform/services/common/nosql/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type VBucketWatcher struct {
	document.DocumentWatcherManager
	key         document.Key
	closeStream chan bool
	logger      *zap.Logger
	cursor      *mongo.ChangeStream
}

func NewVBucketWatcher(key document.Key, collection *mongo.Collection, logger *zap.Logger) *VBucketWatcher {
	result := &VBucketWatcher{
		key:         key,
		logger:      logger,
		closeStream: make(chan bool, 16),
	}
	result.Logger = logger
	err := result.openStream(collection)
	if err != nil {
		return nil
	}
	go result.manageStream()
	return result
}

func (v *VBucketWatcher) openStream(collection *mongo.Collection) error {
	if v.cursor != nil {
		v.cursor.Close(context.TODO())
	}
	changeStreamOpts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	pipeline := mongo.Pipeline{
		{{
			"$match",
			bson.D{
				//{"operationType", "update"},
				{"documentKey._id", v.key.String()},
			},
		}},
	}
	cursor, err := collection.Watch(context.TODO(), pipeline, changeStreamOpts)
	if err != nil {
		v.logger.Warn("Error opening change stream cursor", zap.Error(err))
		return err
	}
	v.cursor = cursor
	return nil
}

func (v *VBucketWatcher) manageStream() {
	defer v.cursor.Close(context.TODO())
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-v.closeStream:
			if !ok {
				v.logger.Info("Closing change stream")
				return
			}
		case <-ticker.C:
			if v.cursor.TryNext(context.TODO()) {
				var changeEvent struct {
					OperationType string `bson:"operationType"`
					//FullDocument      bson.M `bson:"fullDocument"`
					UpdateDescription struct {
						UpdatedFields bson.M `bson:"updatedFields"`
					} `bson:"updateDescription"`
				}
				if err := v.cursor.Decode(&changeEvent); err != nil {
					fmt.Println("decode change event failed:", err)
				}

				//获取变化的数据
				var changed bson.M
				if changeEvent.UpdateDescription.UpdatedFields != nil {
					changed = changeEvent.UpdateDescription.UpdatedFields
				} else {
					fmt.Println("invalid change event", changeEvent.OperationType)
					continue
				}
				data, err := bson.Marshal(changed)
				if err != nil {
					fmt.Println("marshal change event failed:", err)
					continue
				}

				dw := v.WatchersForKey(v.key)
				if len(dw) > 0 {
					v.logger.Info("Mutation",
						zap.String("key", v.key.String()),
						zap.Any("documentKey", v.key.String()),
					)
					for _, w := range dw {
						w.OnDocumentChanged(v.key, data)
					}
				}
			}
			if v.cursor.Err() != nil {
				v.logger.Warn("Error in change stream cursor", zap.Error(v.cursor.Err()))
			}
			if v.cursor.ID() == 0 {
				return
			}

		}
	}
}

func (v *VBucketWatcher) SnapshotMarker(startSeqNo, endSeqNo uint64) {
	v.logger.Info("VBucketWatcher Snapshot",
		zap.String("key", v.key.String()),
		zap.Uint64("startSeqNo", startSeqNo),
		zap.Uint64("endSeqNo", endSeqNo),
	)
}

func (v *VBucketWatcher) Deletion(seqNo, revNo, cas uint64, datatype uint8, vbId uint16, key []byte, value []byte) {
	k := document.NewKeyFromBytesUnchecked(key)
	dw := v.WatchersForKey(k)
	if len(dw) > 0 {
		v.logger.Info("Deletion",
			zap.String("key", v.key.String()),
			zap.String("keys", string(key)),
			zap.Uint64("seqNo", seqNo),
		)
		for _, w := range dw {
			w.OnDocumentDeleted(k)
		}
	}
}

func (v *VBucketWatcher) Expiration(seqNo, revNo, cas uint64, vbId uint16, key []byte) {
	k := document.NewKeyFromBytesUnchecked(key)
	dw := v.WatchersForKey(k)
	if len(dw) > 0 {
		v.logger.Info("Expiration",
			zap.String("key", v.key.String()),
			zap.String("keys", string(key)),
			zap.Uint64("seqNo", seqNo),
		)
		for _, w := range dw {
			w.OnDocumentExpired(k)
		}
	}
}

func (v *VBucketWatcher) Close() {
	v.logger.Info("Close stream end",
		zap.String("key", v.key.String()),
	)
	close(v.closeStream)
}

func (v *VBucketWatcher) StopWatchingDocument(watcher document.DocumentWatcher) bool {
	v.RemoveWatcher(v.key, watcher)
	if len(v.WatchersForKey(v.key)) == 0 {
		v.Close()
		return true
	}
	return false
}
