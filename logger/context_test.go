package logger

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
	ctx := NewContext(nil, NewLogger())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), loggerKey{}, NewLogger())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromContext not works")
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), NewLogger())

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}
