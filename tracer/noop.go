package tracer

import (
	"context"
	"time"

	"go.unistack.org/micro/v3/util/id"
)

var _ Tracer = (*noopTracer)(nil)

type noopTracer struct {
	opts  Options
	spans []Span
}

func (t *noopTracer) Spans() []Span {
	return t.spans
}

func (t *noopTracer) Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	options := NewSpanOptions(opts...)
	span := &noopSpan{
		name:      name,
		ctx:       ctx,
		tracer:    t,
		startTime: time.Now(),
		labels:    options.Labels,
		kind:      options.Kind,
	}
	span.spanID.s, _ = id.New()
	span.traceID.s, _ = id.New()
	if span.ctx == nil {
		span.ctx = context.Background()
	}
	return NewSpanContext(ctx, span), span
}

type noopStringer struct {
	s string
}

func (s noopStringer) String() string {
	return s.s
}

func (t *noopTracer) Flush(ctx context.Context) error {
	return nil
}

func (t *noopTracer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *noopTracer) Name() string {
	return t.opts.Name
}

type noopSpan struct {
	ctx        context.Context
	tracer     Tracer
	name       string
	statusMsg  string
	startTime  time.Time
	finishTime time.Time
	traceID    noopStringer
	spanID     noopStringer
	labels     []interface{}
	kind       SpanKind
	status     SpanStatus
}

func (s *noopSpan) TraceID() string {
	return s.traceID.String()
}

func (s *noopSpan) SpanID() string {
	return s.spanID.String()
}

func (s *noopSpan) Finish(_ ...SpanOption) {
	s.finishTime = time.Now()
}

func (s *noopSpan) Context() context.Context {
	return s.ctx
}

func (s *noopSpan) Tracer() Tracer {
	return s.tracer
}

func (s *noopSpan) AddEvent(name string, opts ...EventOption) {
}

func (s *noopSpan) SetName(name string) {
	s.name = name
}

func (s *noopSpan) AddLogs(kv ...interface{}) {
}

func (s *noopSpan) AddLabels(kv ...interface{}) {
	s.labels = append(s.labels, kv...)
}

func (s *noopSpan) Kind() SpanKind {
	return s.kind
}

func (s *noopSpan) Status() (SpanStatus, string) {
	return s.status, s.statusMsg
}

func (s *noopSpan) SetStatus(st SpanStatus, msg string) {
	s.status = st
	s.statusMsg = msg
}

// NewTracer returns new memory tracer
func NewTracer(opts ...Option) Tracer {
	return &noopTracer{
		opts: NewOptions(opts...),
	}
}
