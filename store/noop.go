package store

import "context"

type NoopStore struct {
	opts Options
}

func (n *NoopStore) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *NoopStore) Options() Options {
	return n.opts
}

func (n *NoopStore) String() string {
	return "noop"
}

func (n *NoopStore) Read(ctx context.Context, key string, opts ...ReadOption) ([]*Record, error) {
	return []*Record{}, nil
}

func (n *NoopStore) Write(ctx context.Context, r *Record, opts ...WriteOption) error {
	return nil
}

func (n *NoopStore) Delete(ctx context.Context, key string, opts ...DeleteOption) error {
	return nil
}

func (n *NoopStore) List(ctx context.Context, opts ...ListOption) ([]string, error) {
	return []string{}, nil
}

func (n *NoopStore) Connect(ctx context.Context) error {
	return nil
}

func (n *NoopStore) Disconnect(ctx context.Context) error {
	return nil
}
