package sql

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/meter"
	"go.unistack.org/micro/v5/tracer"
)

func TestNewOptions_Defaults(t *testing.T) {
	opts := NewOptions()
	if opts.Logger == nil {
		t.Error("expected non-nil logger")
	}
	if opts.Meter == nil {
		t.Error("expected non-nil meter")
	}
	if opts.Tracer == nil {
		t.Error("expected non-nil tracer")
	}
	if opts.MeterStatsInterval != DefaultMeterStatsInterval {
		t.Errorf("expected %v, got %v", DefaultMeterStatsInterval, opts.MeterStatsInterval)
	}
	if opts.LoggerLevel != logger.ErrorLevel {
		t.Errorf("expected ErrorLevel, got %v", opts.LoggerLevel)
	}
}

func TestLogger_Option(t *testing.T) {
	opts := NewOptions(Logger(logger.DefaultLogger))
	if opts.Logger == nil {
		t.Error("expected non-nil logger")
	}
}

func TestMeter_Option(t *testing.T) {
	opts := NewOptions(Meter(meter.DefaultMeter))
	if opts.Meter == nil {
		t.Error("expected non-nil meter")
	}
}

func TestTracer_Option(t *testing.T) {
	opts := NewOptions(Tracer(tracer.DefaultTracer))
	if opts.Tracer == nil {
		t.Error("expected non-nil tracer")
	}
}

func TestMetricInterval_Option(t *testing.T) {
	d := 10 * time.Second
	opts := NewOptions(MetricInterval(d))
	if opts.MeterStatsInterval != d {
		t.Errorf("expected %v, got %v", d, opts.MeterStatsInterval)
	}
}

func TestDatabaseHost_Option(t *testing.T) {
	opts := NewOptions(DatabaseHost("localhost"))
	if opts.DatabaseHost != "localhost" {
		t.Errorf("expected localhost, got %q", opts.DatabaseHost)
	}
}

func TestDatabaseName_Option(t *testing.T) {
	opts := NewOptions(DatabaseName("mydb"))
	if opts.DatabaseName != "mydb" {
		t.Errorf("expected mydb, got %q", opts.DatabaseName)
	}
}

func TestLoggerEnabled_Option(t *testing.T) {
	opts := NewOptions(LoggerEnabled(true))
	if !opts.LoggerEnabled {
		t.Error("expected logger enabled")
	}
}

func TestLoggerLevel_Option(t *testing.T) {
	opts := NewOptions(LoggerLevel(logger.DebugLevel))
	if opts.LoggerLevel != logger.DebugLevel {
		t.Errorf("expected DebugLevel, got %v", opts.LoggerLevel)
	}
}

func TestLoggerObserver_Option(t *testing.T) {
	called := false
	obs := func(ctx context.Context, method, name string, td time.Duration, err error) []any {
		called = true
		return nil
	}
	opts := NewOptions(LoggerObserver(obs))
	opts.LoggerObserver(context.Background(), "m", "n", 0, nil)
	if !called {
		t.Error("expected observer to be called")
	}
}

func TestQueryName(t *testing.T) {
	ctx := QueryName(context.Background(), "my-query")
	if v, ok := ctx.Value(queryNameKey{}).(string); !ok || v != "my-query" {
		t.Errorf("expected 'my-query', got %q", v)
	}
}

func TestQueryName_NilContext(t *testing.T) {
	// nolint: staticcheck
	ctx := QueryName(nil, "q")
	if v, ok := ctx.Value(queryNameKey{}).(string); !ok || v != "q" {
		t.Errorf("expected 'q', got %q", v)
	}
}

func TestGetQueryName_WithValue(t *testing.T) {
	ctx := QueryName(context.Background(), "explicit-query")
	name := getQueryName(ctx)
	if name != "explicit-query" {
		t.Errorf("expected 'explicit-query', got %q", name)
	}
}

func TestGetQueryName_NoValue(t *testing.T) {
	name := getQueryName(context.Background())
	// falls back to getCallerName which returns some caller name
	if name == "" {
		t.Error("expected non-empty caller name")
	}
}

func TestDefaultLoggerObserver(t *testing.T) {
	labels := DefaultLoggerObserver(context.Background(), "Query", "select_users", time.Millisecond, nil)
	if len(labels) == 0 {
		t.Error("expected non-empty labels")
	}
}

func TestDefaultLoggerObserver_WithError(t *testing.T) {
	labels := DefaultLoggerObserver(context.Background(), "Exec", labelUnknown, time.Millisecond, errTest)
	if len(labels) == 0 {
		t.Error("expected non-empty labels")
	}
}

var errTest = &testError{msg: "test error"}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
