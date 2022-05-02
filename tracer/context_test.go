package tracer

import (
	"context"
	"testing"
)

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), tracerKey{}, NewTracer())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromContext not works")
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), NewTracer())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}
