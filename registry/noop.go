package registry

import "fmt"

type noopRegistry struct {
	opts Options
}

func (n *noopRegistry) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *noopRegistry) Options() Options {
	return n.opts
}

func (n *noopRegistry) Register(*Service, ...RegisterOption) error {
	return nil
}

func (n *noopRegistry) Deregister(*Service, ...DeregisterOption) error {
	return nil
}

func (n *noopRegistry) GetService(string, ...GetOption) ([]*Service, error) {
	return []*Service{}, nil
}

func (n *noopRegistry) ListServices(...ListOption) ([]*Service, error) {
	return []*Service{}, nil
}

func (n *noopRegistry) Watch(...WatchOption) (Watcher, error) {
	return nil, fmt.Errorf("not implemented")
}

func (n *noopRegistry) String() string {
	return "noop"
}

// newRegistry returns a new noop registry
func newRegistry(opts ...Option) Registry {
	options := NewOptions()

	for _, o := range opts {
		o(&options)
	}

	return &noopRegistry{opts: options}
}
