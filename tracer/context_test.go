package tracer

import (
	"context"
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
	ctx := NewContext(nil, NewTracer())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

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

func TestSpanFromContext(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	ctx, span := tr.Start(ctx, "test")
	retrieved, ok := SpanFromContext(ctx)
	if !ok {
		t.Fatal("expected span to be in context")
	}
	if retrieved != span {
		t.Error("expected same span")
	}
}

func TestSpanFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	span, ok := SpanFromContext(ctx)
	if ok || span != nil {
		t.Error("expected no span in empty context")
	}
}
