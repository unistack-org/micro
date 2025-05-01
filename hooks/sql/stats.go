package sql

import (
	"context"
	"database/sql"
	"time"
)

type Statser interface {
	Stats() sql.DBStats
}

func NewStatsMeter(ctx context.Context, db Statser, opts ...Option) {
	options := NewOptions(opts...)

	go func() {
		ticker := time.NewTicker(options.MeterStatsInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if db == nil {
					return
				}
				stats := db.Stats()
				options.Meter.Counter(MaxOpenConnections).Set(uint64(stats.MaxOpenConnections))
				options.Meter.Counter(OpenConnections).Set(uint64(stats.OpenConnections))
				options.Meter.Counter(InuseConnections).Set(uint64(stats.InUse))
				options.Meter.Counter(IdleConnections).Set(uint64(stats.Idle))
				options.Meter.Counter(WaitConnections).Set(uint64(stats.WaitCount))
				options.Meter.FloatCounter(BlockedSeconds).Set(stats.WaitDuration.Seconds())
				options.Meter.Counter(MaxIdleClosed).Set(uint64(stats.MaxIdleClosed))
				options.Meter.Counter(MaxIdletimeClosed).Set(uint64(stats.MaxIdleTimeClosed))
				options.Meter.Counter(MaxLifetimeClosed).Set(uint64(stats.MaxLifetimeClosed))
			}
		}
	}()
}
