package meter

import (
	"context"
	"testing"
)

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), meterKey{}, NewMeter())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromContext not works")
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), NewMeter())

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
