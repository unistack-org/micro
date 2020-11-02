package registry

import (
	"context"
	"fmt"
)

type noopRegistry struct {
	opts Options
}

// Init initialize registry
func (n *noopRegistry) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

// Options returns options struct
func (n *noopRegistry) Options() Options {
	return n.opts
}

// Connect opens connection to registry
func (n *noopRegistry) Connect(ctx context.Context) error {
	return nil
}

// Disconnect close connection to registry
func (n *noopRegistry) Disconnect(ctx context.Context) error {
	return nil
}

// Register registers service
func (n *noopRegistry) Register(ctx context.Context, svc *Service, opts ...RegisterOption) error {
	return nil
}

// Deregister deregisters service
func (n *noopRegistry) Deregister(ctx context.Context, svc *Service, opts ...DeregisterOption) error {
	return nil
}

// GetService returns servive info
func (n *noopRegistry) GetService(ctx context.Context, name string, opts ...GetOption) ([]*Service, error) {
	return []*Service{}, nil
}

// ListServices listing services
func (n *noopRegistry) ListServices(ctx context.Context, opts ...ListOption) ([]*Service, error) {
	return []*Service{}, nil
}

// Watch is used to watch for service changes
func (n *noopRegistry) Watch(ctx context.Context, opts ...WatchOption) (Watcher, error) {
	return nil, fmt.Errorf("not implemented")
}

// String returns registry string representation
func (n *noopRegistry) String() string {
	return "noop"
}

// NewRegistry returns a new noop registry
func NewRegistry(opts ...Option) Registry {
	return &noopRegistry{opts: NewOptions(opts...)}
}
