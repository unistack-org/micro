package store

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

func (n *noopStore) Init(opts ...Option) error {
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

func (n *noopStore) Read(key string, opts ...ReadOption) ([]*Record, error) {
	return []*Record{}, nil
}

func (n *noopStore) Write(r *Record, opts ...WriteOption) error {
	return nil
}

func (n *noopStore) Delete(key string, opts ...DeleteOption) error {
	return nil
}

func (n *noopStore) List(opts ...ListOption) ([]string, error) {
	return []string{}, nil
}

func (n *noopStore) Close() error {
	return nil
}
