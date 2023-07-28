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

// NewContext savess broker in context
func NewContext(ctx context.Context, s Broker) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, brokerKey{}, s)
}
