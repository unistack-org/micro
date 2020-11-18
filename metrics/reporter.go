// Package metrics is for instrumentation and debugging
package metrics

import (
	"time"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	DefaultReporter Reporter = NewReporter()
)

// Reporter is an interface for collecting and instrumenting metrics
type Reporter interface {
	Init(...Option) error
	Count(id string, value int64, md metadata.Metadata) error
	Gauge(id string, value float64, md metadata.Metadata) error
	Timing(id string, value time.Duration, md metadata.Metadata) error
	Options() Options
}
