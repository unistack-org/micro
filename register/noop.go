package register

import "context"

type noop struct {
	opts Options
}

func NewRegister(opts ...Option) Register {
	return &noop{
		opts: NewOptions(opts...),
	}
}

func (n *noop) Name() string {
	return n.opts.Name
}

func (n *noop) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *noop) Options() Options {
	return n.opts
}

func (n *noop) Connect(ctx context.Context) error {
	return nil
}

func (n *noop) Disconnect(ctx context.Context) error {
	return nil
}

func (n *noop) Register(ctx context.Context, service *Service, option ...RegisterOption) error {
	return nil
}

func (n *noop) Deregister(ctx context.Context, service *Service, option ...DeregisterOption) error {
	return nil
}

func (n *noop) LookupService(ctx context.Context, s string, option ...LookupOption) ([]*Service, error) {
	return nil, nil
}

func (n *noop) ListServices(ctx context.Context, option ...ListOption) ([]*Service, error) {
	return nil, nil
}

func (n *noop) Watch(ctx context.Context, opts ...WatchOption) (Watcher, error) {
	wOpts := NewWatchOptions(opts...)

	return &watcher{wo: wOpts}, nil
}

func (n *noop) String() string {
	return "noop"
}

type watcher struct {
	wo WatchOptions
}

func (m *watcher) Next() (*Result, error) {
	return nil, nil
}

func (m *watcher) Stop() {}
