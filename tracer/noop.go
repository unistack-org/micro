package tracer

import "context"

type NoopTracer struct{}

func (n *NoopTracer) Init(...Option) error {
	return nil
}

func (n *NoopTracer) Start(ctx context.Context, name string) (context.Context, *Span) {
	return nil, nil
}

func (n *NoopTracer) Finish(*Span) error {
	return nil
}

func (n *NoopTracer) Read(...ReadOption) ([]*Span, error) {
	return nil, nil
}

func NewTracer(opts ...Option) Tracer {
	return &NoopTracer{}
}
