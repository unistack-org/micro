package memory

import (
	"context"
	"time"

	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/tracer"
	"go.unistack.org/micro/v4/util/id"
)

var _ tracer.Tracer = (*Tracer)(nil)

type Tracer struct {
	opts  tracer.Options
	spans []tracer.Span
}

func (t *Tracer) Spans() []tracer.Span {
	return t.spans
}

func (t *Tracer) Start(ctx context.Context, name string, opts ...options.Option) (context.Context, tracer.Span) {
	options := tracer.NewSpanOptions(opts...)
	span := &Span{
		name:      name,
		ctx:       ctx,
		tracer:    t,
		kind:      options.Kind,
		startTime: time.Now(),
	}
	span.spanID.s, _ = id.New()
	span.traceID.s, _ = id.New()
	if span.ctx == nil {
		span.ctx = context.Background()
	}
	t.spans = append(t.spans, span)
	return tracer.NewSpanContext(ctx, span), span
}

func (t *Tracer) Flush(_ context.Context) error {
	return nil
}

func (t *Tracer) Init(opts ...options.Option) error {
	var err error
	for _, o := range opts {
		if err = o(&t.opts); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tracer) Name() string {
	return t.opts.Name
}

type noopStringer struct {
	s string
}

func (s noopStringer) String() string {
	return s.s
}

type Span struct {
	ctx        context.Context
	tracer     tracer.Tracer
	name       string
	statusMsg  string
	startTime  time.Time
	finishTime time.Time
	traceID    noopStringer
	spanID     noopStringer
	events     []*Event
	labels     []interface{}
	logs       []interface{}
	kind       tracer.SpanKind
	status     tracer.SpanStatus
}

func (s *Span) Finish(opts ...options.Option) {
	s.finishTime = time.Now()
}

func (s *Span) Context() context.Context {
	return s.ctx
}

func (s *Span) Tracer() tracer.Tracer {
	return s.tracer
}

type Event struct {
	name   string
	labels []interface{}
}

func (s *Span) AddEvent(name string, opts ...options.Option) {
	options := tracer.NewEventOptions(opts...)
	s.events = append(s.events, &Event{name: name, labels: options.Labels})
}

func (s *Span) SetName(name string) {
	s.name = name
}

func (s *Span) AddLogs(kv ...interface{}) {
	s.logs = append(s.logs, kv...)
}

func (s *Span) AddLabels(kv ...interface{}) {
	s.labels = append(s.labels, kv...)
}

func (s *Span) Kind() tracer.SpanKind {
	return s.kind
}

func (s *Span) TraceID() string {
	return s.traceID.String()
}

func (s *Span) SpanID() string {
	return s.spanID.String()
}

func (s *Span) Status() (tracer.SpanStatus, string) {
	return s.status, s.statusMsg
}

func (s *Span) SetStatus(st tracer.SpanStatus, msg string) {
	s.status = st
	s.statusMsg = msg
}

// NewTracer returns new memory tracer
func NewTracer(opts ...options.Option) *Tracer {
	return &Tracer{
		opts: tracer.NewOptions(opts...),
	}
}
