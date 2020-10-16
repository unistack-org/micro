package registry

import (
	"context"
	"fmt"
)

type NoopRegistry struct {
	opts Options
}

func (n *NoopRegistry) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *NoopRegistry) Options() Options {
	return n.opts
}

func (n *NoopRegistry) Connect(ctx context.Context) error {
	return nil
}

func (n *NoopRegistry) Disconnect(ctx context.Context) error {
	return nil
}

func (n *NoopRegistry) Register(*Service, ...RegisterOption) error {
	return nil
}

func (n *NoopRegistry) Deregister(*Service, ...DeregisterOption) error {
	return nil
}

func (n *NoopRegistry) GetService(string, ...GetOption) ([]*Service, error) {
	return []*Service{}, nil
}

func (n *NoopRegistry) ListServices(...ListOption) ([]*Service, error) {
	return []*Service{}, nil
}

func (n *NoopRegistry) Watch(...WatchOption) (Watcher, error) {
	return nil, fmt.Errorf("not implemented")
}

func (n *NoopRegistry) String() string {
	return "noop"
}

// NewRegistry returns a new noop registry
func NewRegistry(opts ...Option) Registry {
	options := NewOptions(opts...)
	return &NoopRegistry{opts: options}
}
