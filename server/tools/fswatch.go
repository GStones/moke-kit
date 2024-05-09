package tools

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gopkg.in/fsnotify.v1"
)

// Watch watches the filesystem path `path`. When anything changes, changes are
// batched for the period `batchFor`, then `processEvent` is called.
//
// Returns a destroy() function to terminate the watch.
func Watch(logger *zap.Logger, path string, batchFor time.Duration, processEvent func()) (func(), error) {
	logger = logger.With(zap.String("watch files", path))
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	destroy := func() {
		cancel()
		_ = watcher.Close()
	}
	if err := watcher.Add(path); err != nil {
		destroy()
		return nil, err
	}

	go batchWatch(ctx, batchFor, watcher.Events, watcher.Errors, processEvent, func(e error) {
		logger.Error("watcher error", zap.Error(e))
	})
	return destroy, nil
}

// batchWatch: watch for events; when an event occurs, keep draining events for duration `batchFor`, then call processEvent().
// Intended for batching of rapid-fire events where we want to process the batch once, like filesystem update notifications.
func batchWatch(
	ctx context.Context,
	batchFor time.Duration,
	events chan fsnotify.Event,
	errors chan error,
	processEvent func(),
	onError func(error),
) {
	timer := time.NewTimer(0)
	var timerCh <-chan time.Time

	for {
		select {
		// start a timer when an event occurs, otherwise ignore event
		case <-events:
			if timerCh == nil {
				timer.Reset(batchFor)
				timerCh = timer.C
			}

		// on timer, run the batch; nil channels are silently ignored
		case <-timerCh:
			processEvent()
			timerCh = nil

		// handle errors
		case err := <-errors:
			onError(err)

		// on cancel, abort
		case <-ctx.Done():
			return
		}
	}
}
