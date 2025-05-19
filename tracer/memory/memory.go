package memory

import (
	"context"
	"time"

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

func (t *Tracer) Start(ctx context.Context, name string, opts ...tracer.SpanOption) (context.Context, tracer.Span) {
	options := tracer.NewSpanOptions(opts...)
	span := &Span{
		name:      name,
		ctx:       ctx,
		tracer:    t,
		labels:    options.Labels,
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

type memoryStringer struct {
	s string
}

func (s memoryStringer) String() string {
	return s.s
}

func (t *Tracer) Enabled() bool {
	return t.opts.Enabled
}

func (t *Tracer) Flush(_ context.Context) error {
	return nil
}

func (t *Tracer) Init(opts ...tracer.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *Tracer) Name() string {
	return t.opts.Name
}

type Span struct {
	ctx        context.Context
	tracer     tracer.Tracer
	name       string
	statusMsg  string
	startTime  time.Time
	finishTime time.Time
	traceID    memoryStringer
	spanID     memoryStringer
	events     []*Event
	labels     []interface{}
	logs       []interface{}
	kind       tracer.SpanKind
	status     tracer.SpanStatus
}

func (s *Span) Finish(_ ...tracer.SpanOption) {
	s.finishTime = time.Now()
}

func (s *Span) Context() context.Context {
	return s.ctx
}

func (s *Span) Tracer() tracer.Tracer {
	return s.tracer
}

func (s *Span) IsRecording() bool {
	return true
}

type Event struct {
	name   string
	labels []interface{}
}

func (s *Span) AddEvent(name string, opts ...tracer.EventOption) {
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
func NewTracer(opts ...tracer.Option) *Tracer {
	return &Tracer{
		opts: tracer.NewOptions(opts...),
	}
}
