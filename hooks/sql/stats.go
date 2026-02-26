package sql

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

type Statser interface {
	Stats() sql.DBStats
}

func NewStatsMeter(ctx context.Context, db Statser, opts ...Option) {
	if db == nil {
		return
	}

	options := NewOptions(opts...)

	var (
		statsMu                                                     sync.Mutex
		lastUpdated                                                 time.Time
		maxOpenConnections, openConnections, inUse, idle, waitCount float64
		maxIdleClosed, maxIdleTimeClosed, maxLifetimeClosed         float64
		waitDuration                                                float64
	)

	updateFn := func() {
		statsMu.Lock()
		defer statsMu.Unlock()

		if time.Since(lastUpdated) < options.MeterStatsInterval {
			return
		}

		stats := db.Stats()

		maxOpenConnections = float64(stats.MaxOpenConnections)
		openConnections = float64(stats.OpenConnections)
		inUse = float64(stats.InUse)
		idle = float64(stats.Idle)
		waitCount = float64(stats.WaitCount)
		maxIdleClosed = float64(stats.MaxIdleClosed)
		maxIdleTimeClosed = float64(stats.MaxIdleTimeClosed)
		maxLifetimeClosed = float64(stats.MaxLifetimeClosed)
		waitDuration = stats.WaitDuration.Seconds()

		lastUpdated = time.Now()
	}

	options.Meter.Gauge(MaxOpenConnections, func() float64 {
		updateFn()
		return maxOpenConnections
	})

	options.Meter.Gauge(OpenConnections, func() float64 {
		updateFn()
		return openConnections
	})

	options.Meter.Gauge(InuseConnections, func() float64 {
		updateFn()
		return inUse
	})

	options.Meter.Gauge(IdleConnections, func() float64 {
		updateFn()
		return idle
	})

	options.Meter.Gauge(WaitConnections, func() float64 {
		updateFn()
		return waitCount
	})

	options.Meter.Gauge(BlockedSeconds, func() float64 {
		updateFn()
		return waitDuration
	})

	options.Meter.Gauge(MaxIdleClosed, func() float64 {
		updateFn()
		return maxIdleClosed
	})

	options.Meter.Gauge(MaxIdletimeClosed, func() float64 {
		updateFn()
		return maxIdleTimeClosed
	})

	options.Meter.Gauge(MaxLifetimeClosed, func() float64 {
		updateFn()
		return maxLifetimeClosed
	})
}
