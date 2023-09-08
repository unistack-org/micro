package tracer

import (
	"context"
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
		name:   name,
		ctx:    ctx,
		tracer: t,
		labels: options.Labels,
		kind:   options.Kind,
	}
	if span.ctx == nil {
		span.ctx = context.Background()
	}
	t.spans = append(t.spans, span)
	return NewSpanContext(ctx, span), span
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

type noopEvent struct {
	name   string
	labels []interface{}
}

type noopSpan struct {
	ctx       context.Context
	tracer    Tracer
	name      string
	statusMsg string
	events    []*noopEvent
	labels    []interface{}
	logs      []interface{}
	kind      SpanKind
	status    SpanStatus
}

func (s *noopSpan) Finish(opts ...SpanOption) {
}

func (s *noopSpan) Context() context.Context {
	return s.ctx
}

func (s *noopSpan) Tracer() Tracer {
	return s.tracer
}

func (s *noopSpan) AddEvent(name string, opts ...EventOption) {
	options := NewEventOptions(opts...)
	s.events = append(s.events, &noopEvent{name: name, labels: options.Labels})
}

func (s *noopSpan) SetName(name string) {
	s.name = name
}

func (s *noopSpan) AddLogs(kv ...interface{}) {
	s.logs = append(s.logs, kv...)
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
