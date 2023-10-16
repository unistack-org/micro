package tracer

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"go.unistack.org/micro/v4/logger"
)

func TestLoggerWithTracer(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	if err := logger.Init(logger.WithOutput(buf)); err != nil {
		t.Fatal(err)
	}
	var span Span
	ctx, span = DefaultTracer.Start(ctx, "test")
	logger.Info(ctx, "msg")
	if !strings.Contains(buf.String(), span.TraceID()) {
		t.Fatalf("log does not contains tracer id")
	}
}
