package tracer

import "context"

type noopTracer struct {
	opts Options
}

// Init initilize tracer
func (n *noopTracer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

// Start starts new span
func (n *noopTracer) Start(ctx context.Context, name string) (context.Context, *Span) {
	return nil, nil
}

// Finish finishes span
func (n *noopTracer) Finish(*Span) error {
	return nil
}

// Read reads span
func (n *noopTracer) Read(...ReadOption) ([]*Span, error) {
	return nil, nil
}

// NewTracer returns new noop tracer
func NewTracer(opts ...Option) Tracer {
	return &noopTracer{opts: NewOptions(opts...)}
}
