package registry

import (
	"errors"
)

type noopRegistry struct{}

func (n *noopRegistry) Init(...Option) error {
	return nil
}

func (n *noopRegistry) Options() Options {
	return Options{}
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
	return nil, errors.New("not implemented")
}

func (n *noopRegistry) String() string {
	return "noop"
}

// NewRegistry returns a new noop registry
func NewRegistry(...Option) Registry {
	return &noopRegistry{}
}
