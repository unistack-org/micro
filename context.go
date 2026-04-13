package micro

import "context"

type serviceKey struct{}

var serviceKeyVal = serviceKey{}

// FromContext retrieves a Service from the Context.
func FromContext(ctx context.Context) (Service, bool) {
	if ctx == nil {
		return nil, false
	}
	s, ok := ctx.Value(serviceKeyVal).(Service)
	return s, ok
}

// NewContext returns a new Context with the Service embedded within it.
func NewContext(ctx context.Context, s Service) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, serviceKeyVal, s)
}
