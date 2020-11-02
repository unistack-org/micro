// Package transport provides a tunnel transport
package transport

import (
	"context"
	"fmt"

	"github.com/unistack-org/micro/v3/network/transport"
	"github.com/unistack-org/micro/v3/network/tunnel"
)

type tunTransport struct {
	options transport.Options

	tunnel tunnel.Tunnel
}

type tunnelKey struct{}

type transportKey struct{}

func (t *tunTransport) Init(opts ...transport.Option) error {
	for _, o := range opts {
		o(&t.options)
	}

	// close the existing tunnel
	if t.tunnel != nil {
		t.tunnel.Close(context.TODO())
	}

	// get the tunnel
	tun, ok := t.options.Context.Value(tunnelKey{}).(tunnel.Tunnel)
	if !ok {
		return fmt.Errorf("tunnel not set")
	}

	// get the transport
	tr, ok := t.options.Context.Value(transportKey{}).(transport.Transport)
	if ok {
		tun.Init(tunnel.Transport(tr))
	}

	// set the tunnel
	t.tunnel = tun

	return nil
}

func (t *tunTransport) Dial(ctx context.Context, addr string, opts ...transport.DialOption) (transport.Client, error) {
	if err := t.tunnel.Connect(ctx); err != nil {
		return nil, err
	}

	c, err := t.tunnel.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (t *tunTransport) Listen(ctx context.Context, addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	if err := t.tunnel.Connect(ctx); err != nil {
		return nil, err
	}

	l, err := t.tunnel.Listen(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &tunListener{l}, nil
}

func (t *tunTransport) Options() transport.Options {
	return t.options
}

func (t *tunTransport) String() string {
	return "tunnel"
}

// NewTransport honours the initialiser used in
func NewTransport(opts ...transport.Option) transport.Transport {
	t := &tunTransport{
		options: transport.Options{},
	}

	// initialise
	t.Init(opts...)

	return t
}

// WithTunnel sets the internal tunnel
func WithTunnel(t tunnel.Tunnel) transport.Option {
	return func(o *transport.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, tunnelKey{}, t)
	}
}

// WithTransport sets the internal transport
func WithTransport(t transport.Transport) transport.Option {
	return func(o *transport.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, transportKey{}, t)
	}
}
