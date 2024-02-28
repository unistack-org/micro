package register

import (
	"context"
	"go.unistack.org/micro/v4/register/memory"
	"testing"
)

func TestFromNilContext(t *testing.T) {
	// nolint: staticcheck
	c, ok := FromContext(nil)
	if ok || c != nil {
		t.Fatal("FromContext not works")
	}
}

func TestNewNilContext(t *testing.T) {
	// nolint: staticcheck
	ctx := NewContext(nil, memory.NewRegister())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), registerKey{}, memory.NewRegister())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromContext not works")
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), memory.NewRegister())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
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
