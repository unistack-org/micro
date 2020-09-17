package store

import "context"

type noopStore struct {
	opts Options
}

func newStore(opts ...Option) Store {
	options := NewOptions()
	for _, o := range opts {
		o(&options)
	}
	return &noopStore{opts: options}
}

func (n *noopStore) Init(ctx context.Context, opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *noopStore) Options() Options {
	return n.opts
}

func (n *noopStore) String() string {
	return "noop"
}

func (n *noopStore) Read(ctx context.Context, key string, opts ...ReadOption) ([]*Record, error) {
	return []*Record{}, nil
}

func (n *noopStore) Write(ctx context.Context, r *Record, opts ...WriteOption) error {
	return nil
}

func (n *noopStore) Delete(ctx context.Context, key string, opts ...DeleteOption) error {
	return nil
}

func (n *noopStore) List(ctx context.Context, opts ...ListOption) ([]string, error) {
	return []string{}, nil
}

func (n *noopStore) Close(ctx context.Context) error {
	return nil
}
