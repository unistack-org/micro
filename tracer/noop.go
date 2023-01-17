package tracer

import (
	"context"
)

type noopTracer struct {
	opts Options
}

func (t *noopTracer) Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	span := &noopSpan{
		name:   name,
		ctx:    ctx,
		tracer: t,
		opts:   NewSpanOptions(opts...),
	}
	if span.ctx == nil {
		span.ctx = context.Background()
	}
	return NewSpanContext(ctx, span), span
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
	ctx    context.Context
	tracer Tracer
	name   string
	opts   SpanOptions
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
}

func (s *noopSpan) SetName(name string) {
	s.name = name
}

func (s *noopSpan) SetLabels(labels ...interface{}) {
	s.opts.Labels = labels
}

func (s *noopSpan) AddLabels(labels ...interface{}) {
	s.opts.Labels = append(s.opts.Labels, labels...)
}

func (s *noopSpan) Kind() SpanKind {
	return s.opts.Kind
}

// NewTracer returns new memory tracer
func NewTracer(opts ...Option) Tracer {
	return &noopTracer{
		opts: NewOptions(opts...),
	}
}
