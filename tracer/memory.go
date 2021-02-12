package tracer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unistack-org/micro/v3/util/ring"
)

type tracer struct {
	opts Options
	// ring buffer of traces
	buffer *ring.Buffer
}

func (t *tracer) Read(opts ...ReadOption) ([]*Span, error) {
	var options ReadOptions
	for _, o := range opts {
		o(&options)
	}

	sp := t.buffer.Get(t.buffer.Size())

	spans := make([]*Span, 0, len(sp))

	for _, span := range sp {
		val := span.Value.(*Span)
		// skip if trace id is specified and doesn't match
		if len(options.Trace) > 0 && val.Trace != options.Trace {
			continue
		}
		spans = append(spans, val)
	}

	return spans, nil
}

func (t *tracer) Start(ctx context.Context, name string) (context.Context, *Span) {
	span := &Span{
		Name:     name,
		Trace:    uuid.New().String(),
		Id:       uuid.New().String(),
		Started:  time.Now(),
		Metadata: make(map[string]string),
	}

	// return span if no context
	if ctx == nil {
		return NewContext(context.Background(), span.Trace, span.Id), span
	}
	traceID, parentSpanID, ok := FromContext(ctx)
	// If the trace can not be found in the header,
	// that means this is where the trace is created.
	if !ok {
		return NewContext(ctx, span.Trace, span.Id), span
	}

	// set trace id
	span.Trace = traceID
	// set parent
	span.Parent = parentSpanID

	// return the span
	return NewContext(ctx, span.Trace, span.Id), span
}

func (t *tracer) Finish(s *Span) error {
	// set finished time
	s.Duration = time.Since(s.Started)
	// save the span
	t.buffer.Put(s)

	return nil
}

func (t *tracer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *tracer) Lookup(ctx context.Context) (*Span, error) {
	return nil, nil
}

func (t *tracer) Name() string {
	return t.opts.Name
}

func NewTracer(opts ...Option) Tracer {
	return &tracer{
		opts: NewOptions(opts...),
		// the last 256 requests
		buffer: ring.New(256),
	}
}
