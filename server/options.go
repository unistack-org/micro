package server

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/tracer"
	"go.unistack.org/micro/v4/util/id"
)

// Options server struct
type Options struct {
	// Context holds the external options and can be used for server shutdown
	Context context.Context
	// Register holds the register
	Register register.Register
	// Tracer holds the tracer
	Tracer tracer.Tracer
	// Logger holds the logger
	Logger logger.Logger
	// Meter holds the meter
	Meter meter.Meter
	// Listener may be passed if already created
	Listener net.Listener
	// Wait group
	Wait *sync.WaitGroup
	// TLSConfig specifies tls.Config for secure serving
	TLSConfig *tls.Config
	// Metadata holds the server metadata
	Metadata metadata.Metadata
	// RegisterCheck run before register server
	RegisterCheck func(context.Context) error
	// Codecs map to handle content-type
	Codecs map[string]codec.Codec
	// ID holds the id of the server
	ID string
	// Namespace for te server
	Namespace string
	// Name holds the server name
	Name string
	// Address holds the server address
	Address string
	// Advertise holds the advertise address
	Advertise string
	// Version holds the server version
	Version string
	// RegisterAttempts holds the number of register attempts before error
	RegisterAttempts int
	// RegisterInterval holds he interval for re-register
	RegisterInterval time.Duration
	// RegisterTTL specifies TTL for register record
	RegisterTTL time.Duration
	// MaxConn limits number of connections
	MaxConn int
	// DeregisterAttempts holds the number of deregister attempts before error
	DeregisterAttempts int
	// Hooks may contains HandleWrapper or Server func wrapper
	Hooks options.Hooks
	// GracefulTimeout timeout for graceful stop server
	GracefulTimeout time.Duration
}

// NewOptions returns new options struct with default or passed values
func NewOptions(opts ...options.Option) Options {
	options := Options{
		Codecs:           make(map[string]codec.Codec),
		Context:          context.Background(),
		Metadata:         metadata.New(0),
		RegisterInterval: DefaultRegisterInterval,
		RegisterTTL:      DefaultRegisterTTL,
		RegisterCheck:    DefaultRegisterCheck,
		Logger:           logger.DefaultLogger,
		Meter:            meter.DefaultMeter,
		Tracer:           tracer.DefaultTracer,
		Register:         register.DefaultRegister,
		Address:          DefaultAddress,
		Name:             DefaultName,
		Version:          DefaultVersion,
		ID:               id.Must(),
		GracefulTimeout:  DefaultGracefulTimeout,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// ID unique server id
func ID(id string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, id, ".ID")
	}
}

// Version of the service
func Version(v string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, v, ".Version")
	}
}

// Advertise the address to advertise for discovery - host:port
func Advertise(a string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, a, ".Advertise")
	}
}

// RegisterCheck run func before register service
func RegisterCheck(fn func(context.Context) error) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".RegisterCheck")
	}
}

// RegisterTTL registers service with a TTL
func RegisterTTL(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".RegisterTTL")
	}
}

// RegisterInterval registers service with at interval
func RegisterInterval(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".RegisterInterval")
	}
}

// Wait tells the server to wait for requests to finish before exiting
// If `wg` is nil, server only wait for completion of rpc handler.
// For user need finer grained control, pass a concrete `wg` here, server will
// wait against it on stop.
func Wait(wg *sync.WaitGroup) options.Option {
	if wg == nil {
		wg = new(sync.WaitGroup)
	}
	return func(src interface{}) error {
		return options.Set(src, wg, ".Wait")
	}
}

// MaxConn specifies maximum number of max simultaneous connections to server
func MaxConn(n int) options.Option {
	return func(src interface{}) error {
		return options.Set(src, n, ".MaxConn")
	}
}

// Listener specifies the net.Listener to use instead of the default
func Listener(nl net.Listener) options.Option {
	return func(src interface{}) error {
		return options.Set(src, nl, ".Listener")
	}
}

func GracefulTimeout(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".GracefulTimeout")
	}
}

// HandleOptions struct
type HandleOptions struct {
	// Context holds external options
	Context context.Context
	// Metadata for handler
	Metadata map[string]metadata.Metadata
}

// NewHandleOptions creates new HandleOptions
func NewHandleOptions(opts ...options.Option) HandleOptions {
	options := HandleOptions{
		Context:  context.Background(),
		Metadata: make(map[string]metadata.Metadata),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
