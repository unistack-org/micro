package tracer

import (
	"context"
	"testing"
)

func TestNoopTracer_Name(t *testing.T) {
	tr := NewTracer()
	if name := tr.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopTracer_Init(t *testing.T) {
	tr := NewTracer()
	if err := tr.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopTracer_Enabled(t *testing.T) {
	tr := NewTracer()
	if enabled := tr.Enabled(); enabled {
		t.Error("expected tracer to be disabled")
	}
}

func TestNoopTracer_Start(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	ctx, span := tr.Start(ctx, "test-span")
	if span == nil {
		t.Error("expected non-nil span")
	}
	if span.Tracer() != tr {
		t.Error("expected span's tracer to be the same")
	}
}

func TestNoopTracer_Flush(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	if err := tr.Flush(ctx); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMustContext_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from MustContext")
		}
	}()
	MustContext(context.Background())
}

func TestMustContext_Valid(t *testing.T) {
	tr := NewTracer()
	ctx := NewContext(context.Background(), tr)
	retrieved := MustContext(ctx)
	if retrieved != tr {
		t.Error("expected same tracer")
	}
}

func TestSpanMustContext_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from SpanMustContext")
		}
	}()
	SpanMustContext(context.Background())
}

func TestSpanMustContext_Valid(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	ctx, span := tr.Start(ctx, "test")
	retrieved := SpanMustContext(ctx)
	if retrieved != span {
		t.Error("expected same span")
	}
}

func TestSetOption(t *testing.T) {
	key := "test-key"
	val := "test-val"
	opt := SetOption(key, val)
	opts := NewOptions(opt)
	v := opts.Context.Value(key)
	if v != val {
		t.Errorf("expected %v, got %v", val, v)
	}
}

func TestNoopSpan_Methods(t *testing.T) {
	tr := NewTracer()
	ctx := context.Background()
	_, span := tr.Start(ctx, "test")

	// Test TraceID and SpanID (noop uses uuid.Nil)
	if span.TraceID() == "" {
		t.Error("expected non-empty trace ID")
	}
	if span.SpanID() == "" {
		t.Error("expected non-empty span ID")
	}

	// Test IsRecording
	if span.IsRecording() {
		t.Error("expected noop span to not be recording")
	}

	// Test SetName
	span.SetName("new-name")

	// Test AddLabels
	span.AddLabels("k", "v")

	// Test AddEvent
	span.AddEvent("event")

	// Test AddLogs
	span.AddLogs("log-k", "log-v")

	// Test SetStatus
	span.SetStatus(SpanStatusOK, "ok")

	// Test Status
	status, msg := span.Status()
	if status != SpanStatusOK || msg != "ok" {
		t.Error("unexpected status")
	}

	// Test Kind
	span.Kind()

	// Test Context
	span.Context()

	// Test Finish
	span.Finish()
}

func TestSpanStatus_String(t *testing.T) {
	if SpanStatusUnset.String() != "Unset" {
		t.Errorf("unexpected string for SpanStatusUnset")
	}
	if SpanStatusError.String() != "Error" {
		t.Errorf("unexpected string for SpanStatusError")
	}
	if SpanStatusOK.String() != "OK" {
		t.Errorf("unexpected string for SpanStatusOK")
	}
}

func TestSpanKind_String(t *testing.T) {
	if SpanKindUnspecified.String() != "unspecified" {
		t.Errorf("unexpected string for SpanKindUnspecified")
	}
	if SpanKindInternal.String() != "internal" {
		t.Errorf("unexpected string for SpanKindInternal")
	}
	if SpanKindServer.String() != "server" {
		t.Errorf("unexpected string for SpanKindServer")
	}
}

func TestSpanOptions(t *testing.T) {
	// Test WithSpanLabels
	opts := NewSpanOptions(WithSpanLabels("key", "value"))
	if len(opts.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(opts.Labels))
	}

	// Test WithSpanStatus
	opts = NewSpanOptions(WithSpanStatus(SpanStatusError, "err"))
	if opts.Status != SpanStatusError || opts.StatusMsg != "err" {
		t.Error("unexpected span options")
	}

	// Test WithSpanKind
	opts = NewSpanOptions(WithSpanKind(SpanKindClient))
	if opts.Kind != SpanKindClient {
		t.Error("unexpected span kind")
	}

	// Test WithSpanRecord
	opts = NewSpanOptions(WithSpanRecord(false))
	if opts.Record {
		t.Error("expected record to be false")
	}
}

func TestEventOptions(t *testing.T) {
	opts := NewEventOptions(WithEventLabels("k", "v"))
	if len(opts.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(opts.Labels))
	}
}

func TestOptions(t *testing.T) {
	// Test Name option
	tr := NewTracer(Name("test-name"))
	if tr.Name() != "test-name" {
		t.Errorf("expected 'test-name', got %q", tr.Name())
	}

	// Test Disabled option
	tr = NewTracer(Disabled(true))
	if tr.Enabled() {
		t.Error("expected tracer to be disabled")
	}

	// Test Logger option (no-op for noop tracer, but covers the option)
	_ = NewTracer(Logger(nil))
}
