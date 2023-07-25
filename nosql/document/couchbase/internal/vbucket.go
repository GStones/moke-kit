package internal

import (
	"github.com/gstones/platform/services/common/nosql/document"
	"sync"
	"time"

	"github.com/gstones/platform/services/common/utils"
	"go.uber.org/zap"
	"gopkg.in/couchbase/gocbcore.v7"
)

type vBucketWatcher struct {
	document.DocumentWatcherManager

	vbId          uint16
	bucketWatcher *BucketWatcher
	doOpenStream  chan bool
}

func NewVBucketWatcher(vbId uint16, bucketWatcher *BucketWatcher) *vBucketWatcher {
	result := &vBucketWatcher{
		vbId:          vbId,
		bucketWatcher: bucketWatcher,
		doOpenStream:  make(chan bool, 16),
	}

	result.Logger = bucketWatcher.logger

	go result.manageStream()
	result.doOpenStream <- true

	return result
}

func (v *vBucketWatcher) manageStream() {
	bw := v.bucketWatcher
	logger := bw.logger

	client := bw.stream.IoRouter()
	wg := sync.WaitGroup{}

	openStream := func() bool {
		result := false

		var largestSeqNo gocbcore.SeqNo
		var vbUuid gocbcore.VbUuid

		wg.Add(1)

		onOpenComplete := func(entries []gocbcore.FailoverEntry, err error) {
			if err != nil {
				logger.Warn("OpenStream callback failed",
					zap.Error(err),
					zap.Uint16("vbId", v.vbId),
					zap.Uint64("vbUuid", uint64(vbUuid)),
					zap.Uint64("seqNo", uint64(largestSeqNo)),
				)
			} else {
				for _, entry := range entries {
					if entry.SeqNo > largestSeqNo {
						largestSeqNo = entry.SeqNo
						vbUuid = entry.VbUuid
					}
				}
				logger.Info("OpenStream success",
					zap.Uint16("vbId", v.vbId),
					zap.Uint64("vbUuid", uint64(vbUuid)),
					zap.Uint64("seqNo", uint64(largestSeqNo)),
				)
				result = true
			}
			wg.Done()
		}

		logger.Info("OpenStream",
			zap.Uint16("vbId", v.vbId),
			zap.Uint64("vbUuid", uint64(vbUuid)),
			zap.Uint64("seqNo", uint64(largestSeqNo)),
		)

		if _, err := client.OpenStream(
			v.vbId,
			0,
			vbUuid,
			largestSeqNo,
			0xffffffffffffffff,
			0,
			0,
			v,
			onOpenComplete,
		); err != nil {
			logger.Warn("OpenStream failed",
				zap.Error(err),
				zap.Uint16("vbId", v.vbId),
				zap.Uint64("vbUuid", uint64(vbUuid)),
				zap.Uint64("seqNo", uint64(largestSeqNo)),
			)
			wg.Done()
		}

		wg.Wait()

		if result {
			v.OnAllDocumentsChanged()
		}

		return result
	}

	bo := utils.NewExponentialSleep(100*time.Millisecond, 5*time.Second, 0.5)

	for {
		if doOpen, ok := <-v.doOpenStream; doOpen && ok {
			bo.Reset()

			for !openStream() {
				bo.Sleep()
			}
		} else {
			break
		}
	}
}

func (v *vBucketWatcher) SnapshotMarker(startSeqNo, endSeqNo uint64, vbId uint16, snapshotType gocbcore.SnapshotState) {
	v.bucketWatcher.logger.Info("VBucketWatcher Snapshot",
		zap.Uint16("vbId", v.vbId),
		zap.Uint64("startSeqNo", startSeqNo),
		zap.Uint64("endSeqNo", endSeqNo),
	)
}

func (v *vBucketWatcher) Mutation(
	seqNo, revNo uint64,
	flags, expiry, lockTime uint32,
	cas uint64,
	datatype uint8,
	vbId uint16,
	key []byte,
	value []byte,
) {
	k := document.NewKeyFromBytesUnchecked(key)
	dw := v.WatchersForKey(k)
	if len(dw) > 0 {
		v.bucketWatcher.logger.Info("Mutation",
			zap.Uint16("vbId", v.vbId),
			zap.String("keys", string(key)),
			zap.Uint64("seqNo", seqNo),
		)
		for _, w := range dw {
			w.OnDocumentChanged(k, nil)
		}
	}
}

func (v *vBucketWatcher) Deletion(seqNo, revNo, cas uint64, datatype uint8, vbId uint16, key []byte, value []byte) {
	k := document.NewKeyFromBytesUnchecked(key)
	dw := v.WatchersForKey(k)
	if len(dw) > 0 {
		v.bucketWatcher.logger.Info("Deletion",
			zap.Uint16("vbId", v.vbId),
			zap.String("keys", string(key)),
			zap.Uint64("seqNo", seqNo),
		)
		for _, w := range dw {
			w.OnDocumentDeleted(k)
		}
	}
}

func (v *vBucketWatcher) Expiration(seqNo, revNo, cas uint64, vbId uint16, key []byte) {
	k := document.NewKeyFromBytesUnchecked(key)
	dw := v.WatchersForKey(k)
	if len(dw) > 0 {
		v.bucketWatcher.logger.Info("Expiration",
			zap.Uint16("vbId", v.vbId),
			zap.String("keys", string(key)),
			zap.Uint64("seqNo", seqNo),
		)
		for _, w := range dw {
			w.OnDocumentExpired(k)
		}
	}
}

func (v *vBucketWatcher) End(vbId uint16, err error) {
	v.bucketWatcher.logger.Error("VBucketWatcher stream end",
		zap.Error(err),
		zap.Uint16("vbId", v.vbId),
	)
	v.doOpenStream <- true
}
