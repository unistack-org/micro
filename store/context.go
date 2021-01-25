package store

import (
	"context"
)

type storeKey struct{}

// FromContext get store from context
func FromContext(ctx context.Context) (Store, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(storeKey{}).(Store)
	return c, ok
}

// NewContext put store in context
func NewContext(ctx context.Context, c Store) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, storeKey{}, c)
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

// SetReadOption returns a function to setup a context with given value
func SetReadOption(k, v interface{}) ReadOption {
	return func(o *ReadOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetWriteOption returns a function to setup a context with given value
func SetWriteOption(k, v interface{}) WriteOption {
	return func(o *WriteOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetListOption returns a function to setup a context with given value
func SetListOption(k, v interface{}) ListOption {
	return func(o *ListOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetDeleteOption returns a function to setup a context with given value
func SetDeleteOption(k, v interface{}) DeleteOption {
	return func(o *DeleteOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetExistsOption returns a function to setup a context with given value
func SetExistsOption(k, v interface{}) ExistsOption {
	return func(o *ExistsOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
