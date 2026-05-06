package broker

import (
	"context"
	"testing"
)

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), brokerKey{}, NewBroker())
	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromContext not works")
	}
}

func TestFromNilContext(t *testing.T) {
	// nolint: staticcheck
	c, ok := FromContext(nil)
	if ok || c != nil {
		t.Fatal("FromContext not works")
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), NewBroker())
	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

func TestNewNilContext(t *testing.T) {
	// nolint: staticcheck
	ctx := NewContext(nil, NewBroker())
	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

func TestSetSubscribeOption(t *testing.T) {
	type key struct{}
	o := SetSubscribeOption(key{}, "test")
	opts := &SubscribeOptions{}
	o(opts)

	if v, ok := opts.Context.Value(key{}).(string); !ok || v == "" {
		t.Fatal("SetSubscribeOption not works")
	}
}

func TestSetOption(t *testing.T) {
	type key struct{}
	o := SetOption(key{}, "test")
	opts := &Options{}
	o(opts)

	if v, ok := opts.Context.Value(key{}).(string); !ok || v == "" {
		t.Fatal("SetOption not works")
	}
}

func TestWithBroker(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()
	ctx = NewContext(ctx, b)
	retrieved, ok := FromContext(ctx)
	if !ok || retrieved != b {
		t.Error("expected broker to be retrieved from context")
	}
}

func TestBrokerFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	retrieved, ok := FromContext(ctx)
	if ok || retrieved != nil {
		t.Error("expected nil broker from empty context")
	}
}
