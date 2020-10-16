// Package registry is an interface for service discovery
package registry

import (
	"context"
	"errors"
)

const (
	// WildcardDomain indicates any domain
	WildcardDomain = "*"
	// DefaultDomain to use if none was provided in options
	DefaultDomain = "micro"
)

var (
	DefaultRegistry Registry = NewRegistry()
	// ErrNotFound returned when GetService is called and no services found
	ErrNotFound = errors.New("service not found")
	// ErrWatcherStopped returned when when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

// The registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Options() Options
	Connect(context.Context) error
	Disconnect(context.Context) error
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

// Service holds service registry info
type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

// Node holds node registry info
type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

// Endpoint holds endpoint registry info
type Endpoint struct {
	Name     string            `json:"name"`
	Request  *Value            `json:"request"`
	Response *Value            `json:"response"`
	Metadata map[string]string `json:"metadata"`
}

// Valud holds additional kv stuff
type Value struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []*Value `json:"values"`
}

// Option func signature
type Option func(*Options)

// RegisterOption option is used to register service
type RegisterOption func(*RegisterOptions)

// WatchOption option is used to watch service changes
type WatchOption func(*WatchOptions)

// DeregisterOption option is used to deregister service
type DeregisterOption func(*DeregisterOptions)

// GetOption option is used to get service
type GetOption func(*GetOptions)

// ListOption option is used to list services
type ListOption func(*ListOptions)
