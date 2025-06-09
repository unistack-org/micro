package broker

import (
	"context"
)

type brokerKey struct{}

// FromContext returns broker from passed context
func FromContext(ctx context.Context) (Broker, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(brokerKey{}).(Broker)
	return c, ok
}

// MustContext returns broker from passed context
func MustContext(ctx context.Context) Broker {
	b, ok := FromContext(ctx)
	if !ok {
		panic("missing broker")
	}
	return b
}

// NewContext savess broker in context
func NewContext(ctx context.Context, s Broker) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, brokerKey{}, s)
}

// SetSubscribeOption returns a function to setup a context with given value
func SetSubscribeOption(k, v interface{}) SubscribeOption {
	return func(o *SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetMessageOption returns a function to setup a context with given value
func SetMessageOption(k, v interface{}) MessageOption {
	return func(o *MessageOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetOption returns a function to setup a context with given value
func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
