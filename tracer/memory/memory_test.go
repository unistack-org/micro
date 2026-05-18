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

func TestMemoryTracer_Spans(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	if spans := tr.Spans(); len(spans) != 0 {
		t.Errorf("expected empty spans, got %d", len(spans))
	}
	tr.Start(ctx, "span-1")
	tr.Start(ctx, "span-2")
	if spans := tr.Spans(); len(spans) != 2 {
		t.Errorf("expected 2 spans, got %d", len(spans))
	}
}

func TestMemoryTracer_Enabled(t *testing.T) {
	// Disabled(false) means enabled=true (Enabled = !b)
	tr := NewTracer(tracer.Disabled(false))
	if !tr.Enabled() {
		t.Error("expected tracer to be enabled when Disabled(false)")
	}
	// Disabled(true) means enabled=false
	tr2 := NewTracer(tracer.Disabled(true))
	if tr2.Enabled() {
		t.Error("expected tracer to be disabled when Disabled(true)")
	}
}

func TestMemoryTracer_Flush(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	if err := tr.Flush(ctx); err != nil {
		t.Errorf("Flush returned unexpected error: %v", err)
	}
}

func TestMemoryTracer_Init(t *testing.T) {
	tr := NewTracer()
	if err := tr.Init(tracer.Name("custom")); err != nil {
		t.Errorf("Init returned unexpected error: %v", err)
	}
	if tr.Name() != "custom" {
		t.Errorf("expected name 'custom', got %q", tr.Name())
	}
}

func TestSpan_Context(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	_, span := tr.Start(ctx, "test-span")
	if span.Context() == nil {
		t.Error("expected non-nil span context")
	}
}

func TestSpan_Tracer(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	_, span := tr.Start(ctx, "test-span")
	if span.Tracer() != tr {
		t.Error("expected span tracer to be the same tracer")
	}
}

func TestSpan_Kind(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	_, span := tr.Start(ctx, "test-span", tracer.WithSpanKind(tracer.SpanKindClient))
	if span.Kind() != tracer.SpanKindClient {
		t.Errorf("expected SpanKindClient, got %v", span.Kind())
	}
}

func TestSpan_Status(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	_, span := tr.Start(ctx, "test-span")
	span.SetStatus(tracer.SpanStatusError, "something failed")
	st, msg := span.Status()
	if st != tracer.SpanStatusError {
		t.Errorf("expected SpanStatusError, got %v", st)
	}
	if msg != "something failed" {
		t.Errorf("expected 'something failed', got %q", msg)
	}
}

func TestMemoryTracer_StartNilContext(t *testing.T) {
	tr := NewTracer()
	_, span := tr.Start(nil, "nil-ctx-span") //nolint:staticcheck
	if span == nil {
		t.Fatal("expected non-nil span even with nil context")
	}
	if span.Context() == nil {
		t.Error("expected context to be set to Background when nil is passed")
	}
}
