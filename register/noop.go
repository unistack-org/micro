package register

import (
	"context"
)

type noopRegister struct {
	opts Options
}

func (n *noopRegister) Name() string {
	return n.opts.Name
}

// Init initialize register
func (n *noopRegister) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

// Options returns options struct
func (n *noopRegister) Options() Options {
	return n.opts
}

// Connect opens connection to register
func (n *noopRegister) Connect(ctx context.Context) error {
	return nil
}

// Disconnect close connection to register
func (n *noopRegister) Disconnect(ctx context.Context) error {
	return nil
}

// Register registers service
func (n *noopRegister) Register(ctx context.Context, svc *Service, opts ...RegisterOption) error {
	return nil
}

// Deregister deregisters service
func (n *noopRegister) Deregister(ctx context.Context, svc *Service, opts ...DeregisterOption) error {
	return nil
}

// LookupService returns servive info
func (n *noopRegister) LookupService(ctx context.Context, name string, opts ...LookupOption) ([]*Service, error) {
	return []*Service{}, nil
}

// ListServices listing services
func (n *noopRegister) ListServices(ctx context.Context, opts ...ListOption) ([]*Service, error) {
	return []*Service{}, nil
}

// Watch is used to watch for service changes
func (n *noopRegister) Watch(ctx context.Context, opts ...WatchOption) (Watcher, error) {
	return &noopWatcher{done: make(chan struct{}), opts: NewWatchOptions(opts...)}, nil
}

// String returns register string representation
func (n *noopRegister) String() string {
	return "noop"
}

type noopWatcher struct {
	opts WatchOptions
	done chan struct{}
}

func (n *noopWatcher) Next() (*Result, error) {
	<-n.done
	return nil, ErrWatcherStopped
}

func (n *noopWatcher) Stop() {
	close(n.done)
}

// NewRegister returns a new noop register
func NewRegister(opts ...Option) Register {
	return &noopRegister{opts: NewOptions(opts...)}
}
