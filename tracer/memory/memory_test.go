package memory

import (
	"context"
	"testing"

	"go.unistack.org/micro/v5/tracer"
)

func TestMemoryTracer_Name(t *testing.T) {
	tr := NewTracer()
	if name := tr.Name(); name != "memory" {
		t.Errorf("expected 'memory', got %q", name)
	}
}

func TestMemoryTracer_StartFinish(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	_, span := tr.Start(ctx, "test-span")
	if span == nil {
		t.Fatal("expected non-nil span")
	}
	span.Finish()
}

func TestMemoryTracer_SpanOperations(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	_, span := tr.Start(ctx, "test-span")
	span.SetName("new-name")
	span.SetStatus(tracer.SpanStatusOK, "ok")
	span.AddLabels("key", "value")
	span.AddEvent("event")
	span.AddLogs("log-key", "log-value")
	if span.TraceID() == "" {
		t.Error("expected non-empty trace ID")
	}
	if span.SpanID() == "" {
		t.Error("expected non-empty span ID")
	}
	if !span.IsRecording() {
		t.Error("expected span to be recording")
	}
	span.Finish()
}
