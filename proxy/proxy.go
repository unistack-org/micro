// Package proxy is a transparent proxy built on the micro/server
package proxy

import (
	"context"

	"github.com/unistack-org/micro/v3/server"
)

// Proxy can be used as a proxy server for micro services
type Proxy interface {
	// ProcessMessage handles inbound messages
	ProcessMessage(context.Context, server.Message) error
	// ServeRequest handles inbound requests
	ServeRequest(context.Context, server.Request, server.Response) error
	// Name of the proxy protocol
	String() string
}

var (
	DefaultEndpoint = "localhost:9090"
)
