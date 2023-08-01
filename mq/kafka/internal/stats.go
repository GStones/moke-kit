package internal

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// Note: An alias is used instead of a declaration because uconfig parsing fails
// when the underlying time.Duration type is not exposed.
type StatsPeriod = time.Duration

type StatsLogger struct {
	logger *zap.Logger
	period StatsPeriod
	stats  StatsProvider
	done   chan struct{}

	// mutable fields synchronized on the mutex
	mutex   sync.Mutex
	started bool
}

func NewStatsLogger(logger *zap.Logger, period StatsPeriod, stats StatsProvider) *StatsLogger {
	return &StatsLogger{
		logger: logger,
		period: period,
		stats:  stats,
		done:   make(chan struct{}),
	}
}

func (l *StatsLogger) Start() {
	// safely check the started flag
	l.mutex.Lock()
	if l.started {
		l.mutex.Unlock()
		return
	} else {
		l.started = true
		l.mutex.Unlock()
	}

	// kick off a goroutine to periodically report stats in the background
	ticker := time.NewTicker(time.Duration(l.period))
	go func() {
		select {
		case <-ticker.C:
			l.Report()
		case <-l.done:
			ticker.Stop()
			return
		}
	}()
}

func (l *StatsLogger) Stop() {
	l.mutex.Lock()
	l.started = false
	close(l.done)
	l.mutex.Unlock()
}

func (l *StatsLogger) Report() {
	msg, stats := l.stats.Stats()
	field := zap.Any(msg, stats)
	l.logger.Info("stats", field)
}
