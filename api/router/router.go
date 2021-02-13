// Package router provides api service routing
package router

import (
	"net/http"

	"github.com/unistack-org/micro/v3/api"
)

var (
	DefaultRouter Router
)

// Router is used to determine an endpoint for a request
type Router interface {
	// Returns options
	Options() Options
	// Init initialize router
	Init(...Option) error
	// Stop the router
	Close() error
	// Endpoint returns an api.Service endpoint or an error if it does not exist
	Endpoint(r *http.Request) (*api.Service, error)
	// Register endpoint in router
	Register(ep *api.Endpoint) error
	// Deregister endpoint from router
	Deregister(ep *api.Endpoint) error
	// Route returns an api.Service route
	Route(r *http.Request) (*api.Service, error)
	// String representation of router
	String() string
}
