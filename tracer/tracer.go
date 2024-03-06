// Package tracer provides an interface for distributed tracing
package tracer // import "go.unistack.org/micro/v3/tracer"

import (
	"context"
	"fmt"
	"sort"

	"go.unistack.org/micro/v3/logger"
)

// DefaultTracer is the global default tracer
var DefaultTracer Tracer = NewTracer()

var (
	// TraceIDKey is the key used for the trace id in the log call
	TraceIDKey = "trace-id"
	// SpanIDKey is the key used for the span id in the log call
	SpanIDKey = "span-id"
)

func init() {
	logger.DefaultContextAttrFuncs = append(logger.DefaultContextAttrFuncs,
		func(ctx context.Context) []interface{} {
			if span, ok := SpanFromContext(ctx); ok {
				return []interface{}{
					TraceIDKey, span.TraceID(),
					SpanIDKey, span.SpanID(),
				}
			}
			return nil
		})
}

// Tracer is an interface for distributed tracing
type Tracer interface {
	// Name return tracer name
	Name() string
	// Init tracer with options
	Init(...Option) error
	// Start a trace
	Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span)
	// Flush flushes spans
	Flush(ctx context.Context) error
}

type Span interface {
	// Tracer return underlining tracer
	Tracer() Tracer
	// Finish complete and send span
	Finish(opts ...SpanOption)
	// Context return context with span
	Context() context.Context
	// SetName set the span name
	SetName(name string)
	// SetStatus set the span status code and msg
	SetStatus(status SpanStatus, msg string)
	// Status returns span status and msg
	Status() (SpanStatus, string)
	// AddLabels append labels to span
	AddLabels(kv ...interface{})
	// AddEvent append event to span
	AddEvent(name string, opts ...EventOption)
	// AddEvent append event to span
	AddLogs(kv ...interface{})
	// Kind returns span kind
	Kind() SpanKind
	// TraceID returns trace id
	TraceID() string
	// SpanID returns span id
	SpanID() string
}

// sort labels alphabeticaly by label name
type byKey []interface{}

func (k byKey) Len() int           { return len(k) / 2 }
func (k byKey) Less(i, j int) bool { return fmt.Sprintf("%s", k[i*2]) < fmt.Sprintf("%s", k[j*2]) }
func (k byKey) Swap(i, j int) {
	k[i*2], k[j*2] = k[j*2], k[i*2]
	k[i*2+1], k[j*2+1] = k[j*2+1], k[i*2+1]
}

func UniqLabels(labels []interface{}) []interface{} {
	if len(labels)%2 == 1 {
		labels = labels[:len(labels)-1]
	}

	if len(labels) > 2 {
		sort.Sort(byKey(labels))

		idx := 0
		for {
			if labels[idx] == labels[idx+2] {
				copy(labels[idx:], labels[idx+2:])
				labels = labels[:len(labels)-2]
			} else {
				idx += 2
			}
			if idx+2 >= len(labels) {
				break
			}
		}
	}
	return labels
}
