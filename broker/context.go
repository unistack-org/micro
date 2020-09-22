package broker

import (
	"context"
)

type brokerKey struct{}

func FromContext(ctx context.Context) (Broker, bool) {
	c, ok := ctx.Value(brokerKey{}).(Broker)
	return c, ok
}

func NewContext(ctx context.Context, s Broker) context.Context {
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

// SetOption returns a function to setup a context with given value
func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetPublishOption returns a function to setup a context with given value
func SetPublishOption(k, v interface{}) PublishOption {
	return func(o *PublishOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
