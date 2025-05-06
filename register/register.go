// Package register is an interface for service discovery
package register

import (
	"context"
	"errors"

	"go.unistack.org/micro/v4/metadata"
)

const (
	// WildcardNamespace indicates any Namespace
	WildcardNamespace = "*"
)

// DefaultNamespace to use if none was provided in options
var DefaultNamespace = "micro"

var (
	// DefaultRegister is the global default register
	DefaultRegister = NewRegister()
	// ErrNotFound returned when LookupService is called and no services found
	ErrNotFound = errors.New("service not found")
	// ErrWatcherStopped returned when when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

// Register provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Register interface {
	// Name returns register name
	Name() string
	// Init initialize register
	Init(...Option) error
	// Options returns options for register
	Options() Options
	// Connect initialize connect to register
	Connect(context.Context) error
	// Disconnect initialize discconection from register
	Disconnect(context.Context) error
	// Register service in registry
	Register(context.Context, *Service, ...RegisterOption) error
	// Deregister service from registry
	Deregister(context.Context, *Service, ...DeregisterOption) error
	// LookupService in registry
	LookupService(context.Context, string, ...LookupOption) ([]*Service, error)
	// ListServices in registry
	ListServices(context.Context, ...ListOption) ([]*Service, error)
	// Watch registry events
	Watch(context.Context, ...WatchOption) (Watcher, error)
	// String returns registry string representation
	String() string
	// Live returns register liveness
	// Live() bool
	// Ready returns register readiness
	// Ready() bool
}

// Service holds service register info
type Service struct {
	Name      string  `json:"name,omitempty"`
	Version   string  `json:"version,omitempty"`
	Nodes     []*Node `json:"nodes,omitempty"`
	Namespace string  `json:"namespace,omitempty"`
}

// Node holds node register info
type Node struct {
	Metadata metadata.Metadata `json:"metadata,omitempty"`
	ID       string            `json:"id,omitempty"`
	// Address also prefixed with scheme like grpc://xx.xx.xx.xx:1234
	Address string `json:"address,omitempty"`
}

// Option func signature
type Option func(*Options)

// RegisterOption option is used to register service
type RegisterOption func(*RegisterOptions) // nolint: golint,revive

// WatchOption option is used to watch service changes
type WatchOption func(*WatchOptions)

// DeregisterOption option is used to deregister service
type DeregisterOption func(*DeregisterOptions)

// LookupOption option is used to get service
type LookupOption func(*LookupOptions)

// ListOption option is used to list services
type ListOption func(*ListOptions)
