package logger

import (
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.TODO()
	l := NewLogger(WithLevel(TraceLevel))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Trace(ctx, "trace_msg1")
	l.Warn(ctx, "warn_msg1")
	l.Fields(map[string]interface{}{"error": "test"}).Info(ctx, "error message")
	l.Warn(ctx, "first", " ", "second")
}
