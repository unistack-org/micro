package tracer_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/logger/slog"
	"go.unistack.org/micro/v4/tracer"
)

func TestLoggerWithTracer(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	logger.DefaultLogger = slog.NewLogger(logger.WithOutput(buf))

	if err := logger.Init(); err != nil {
		t.Fatal(err)
	}
	var span tracer.Span
	ctx, span = tracer.DefaultTracer.Start(ctx, "test")
	logger.Info(ctx, "msg")
	if !strings.Contains(buf.String(), span.TraceID()) {
		t.Fatalf("log does not contains tracer id")
	}
}
