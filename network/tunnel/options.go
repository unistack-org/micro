package tunnel

import (
	"time"

	"github.com/google/uuid"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/network/transport"
	"github.com/unistack-org/micro/v3/tracer"
)

var (
	// DefaultAddress is default tunnel bind address
	DefaultAddress = ":0"
	// The shared default token
	DefaultToken = "go.micro.tunnel"
)

// Option func
type Option func(*Options)

// Options provides network configuration options
type Options struct {
	Name string
	// Id is tunnel id
	Id string
	// Address is tunnel address
	Address string
	// Nodes are remote nodes
	Nodes []string
	// The shared auth token
	Token string
	// Transport listens to incoming connections
	Transport transport.Transport
	// Logger
	Logger logger.Logger
	// Meter
	Meter meter.Meter
	// Tracer
	Tracer tracer.Tracer
}

// DialOption func
type DialOption func(*DialOptions)

// DialOptions provides dial options
type DialOptions struct {
	// Link specifies the link to use
	Link string
	// specify mode of the session
	Mode Mode
	// Wait for connection to be accepted
	Wait bool
	// the dial timeout
	Timeout time.Duration
}

// ListenOption func
type ListenOption func(*ListenOptions)

// ListenOptions provides listen options
type ListenOptions struct {
	// specify mode of the session
	Mode Mode
	// The read timeout
	Timeout time.Duration
}

// Id sets the tunnel id
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Meter sets the meter
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Address sets the tunnel address
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Nodes specify remote network nodes
func Nodes(n ...string) Option {
	return func(o *Options) {
		o.Nodes = n
	}
}

// Token sets the shared token for auth
func Token(t string) Option {
	return func(o *Options) {
		o.Token = t
	}
}

// Transport listens for incoming connections
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// ListenMode option
func ListenMode(m Mode) ListenOption {
	return func(o *ListenOptions) {
		o.Mode = m
	}
}

// ListenTimeout for reads and writes on the listener session
func ListenTimeout(t time.Duration) ListenOption {
	return func(o *ListenOptions) {
		o.Timeout = t
	}
}

// DialMode multicast sets the multicast option to send only to those mapped
func DialMode(m Mode) DialOption {
	return func(o *DialOptions) {
		o.Mode = m
	}
}

// DialTimeout sets the dial timeout of the connection
func DialTimeout(t time.Duration) DialOption {
	return func(o *DialOptions) {
		o.Timeout = t
	}
}

// DialLink specifies the link to pin this connection to.
// This is not applicable if the multicast option is set.
func DialLink(id string) DialOption {
	return func(o *DialOptions) {
		o.Link = id
	}
}

// DialWait specifies whether to wait for the connection
// to be accepted before returning the session
func DialWait(b bool) DialOption {
	return func(o *DialOptions) {
		o.Wait = b
	}
}

// DefaultOptions returns router default options
func NewOptions(opts ...Option) Options {
	options := Options{
		Id:      uuid.New().String(),
		Address: DefaultAddress,
		Token:   DefaultToken,
		Logger:  logger.DefaultLogger,
		Meter:   meter.DefaultMeter,
		Tracer:  tracer.DefaultTracer,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Tracer to be used for tracing
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
