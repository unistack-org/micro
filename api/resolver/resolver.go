// Package resolver resolves a http request to an endpoint
package resolver

import (
	"errors"
	"net/http"
)

var (
	// ErrNotFound returned when endpoint is not found
	ErrNotFound = errors.New("not found")
	// ErrInvalidPath returned on invalid path
	ErrInvalidPath = errors.New("invalid path")
)

// Resolver resolves requests to endpoints
type Resolver interface {
	Resolve(r *http.Request, opts ...ResolveOption) (*Endpoint, error)
	String() string
}

// Endpoint is the endpoint for a http request
type Endpoint struct {
	// Endpoint name e.g greeter
	Name string
	// HTTP Host e.g example.com
	Host string
	// HTTP Methods e.g GET, POST
	Method string
	// HTTP Path e.g /greeter.
	Path string
	// Domain endpoint exists within
	Domain string
}
