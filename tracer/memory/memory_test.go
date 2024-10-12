package memory

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/logger/slog"
	"go.unistack.org/micro/v3/tracer"
)

func TestLoggerWithTracer(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	logger.DefaultLogger = slog.NewLogger(logger.WithOutput(buf))

	if err := logger.DefaultLogger.Init(); err != nil {
		t.Fatal(err)
	}
	var span tracer.Span
	tr := NewTracer()
	ctx, span = tr.Start(ctx, "test1")

	logger.DefaultLogger.Error(ctx, "my test error", fmt.Errorf("error"))

	if !strings.Contains(buf.String(), span.TraceID()) {
		t.Fatalf("log does not contains trace id: %s", buf.Bytes())
	}

	_, _ = tr.Start(ctx, "test2")

	for _, s := range tr.Spans() {
		_ = s
	}
}
