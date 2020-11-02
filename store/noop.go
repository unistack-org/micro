package store

import "context"

type noopStore struct {
	opts Options
}

func NewStore(opts ...Option) Store {
	return &noopStore{opts: NewOptions(opts...)}
}

// Init initialize store
func (n *noopStore) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

// Options returns Options struct
func (n *noopStore) Options() Options {
	return n.opts
}

// String returns string representation
func (n *noopStore) String() string {
	return "noop"
}

// Read reads store value by key
func (n *noopStore) Read(ctx context.Context, key string, opts ...ReadOption) ([]*Record, error) {
	return []*Record{}, nil
}

// Write writes store record
func (n *noopStore) Write(ctx context.Context, r *Record, opts ...WriteOption) error {
	return nil
}

// Delete removes store value by key
func (n *noopStore) Delete(ctx context.Context, key string, opts ...DeleteOption) error {
	return nil
}

// List lists store
func (n *noopStore) List(ctx context.Context, opts ...ListOption) ([]string, error) {
	return []string{}, nil
}

// Connect connects to store
func (n *noopStore) Connect(ctx context.Context) error {
	return nil
}

// Disconnect disconnects from store
func (n *noopStore) Disconnect(ctx context.Context) error {
	return nil
}
