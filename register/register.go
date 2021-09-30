// Package register is an interface for service discovery
package register

import (
	"context"
	"errors"

	"github.com/unistack-org/micro/v3/metadata"
)

const (
	// WildcardDomain indicates any domain
	WildcardDomain = "*"
	// DefaultDomain to use if none was provided in options
	DefaultDomain = "micro"
)

var (
	// DefaultRegister is the global default register
	DefaultRegister Register = NewRegister()
	// ErrNotFound returned when LookupService is called and no services found
	ErrNotFound = errors.New("service not found")
	// ErrWatcherStopped returned when when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

// Register provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Register interface {
	Name() string
	Init(...Option) error
	Options() Options
	Connect(context.Context) error
	Disconnect(context.Context) error
	Register(context.Context, *Service, ...RegisterOption) error
	Deregister(context.Context, *Service, ...DeregisterOption) error
	LookupService(context.Context, string, ...LookupOption) ([]*Service, error)
	ListServices(context.Context, ...ListOption) ([]*Service, error)
	Watch(context.Context, ...WatchOption) (Watcher, error)
	String() string
}

// Service holds service register info
type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  metadata.Metadata `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

// Node holds node register info
type Node struct {
	Metadata metadata.Metadata `json:"metadata"`
	ID       string            `json:"id"`
	Address  string            `json:"address"`
}

// Endpoint holds endpoint register info
type Endpoint struct {
	Request  string            `json:"request"`
	Response string            `json:"response"`
	Metadata metadata.Metadata `json:"metadata"`
	Name     string            `json:"name"`
}

// Option func signature
type Option func(*Options)

// nolint: golint,revive
// RegisterOption option is used to register service
type RegisterOption func(*RegisterOptions)

// WatchOption option is used to watch service changes
type WatchOption func(*WatchOptions)

// DeregisterOption option is used to deregister service
type DeregisterOption func(*DeregisterOptions)

// LookupOption option is used to get service
type LookupOption func(*LookupOptions)

// ListOption option is used to list services
type ListOption func(*ListOptions)
