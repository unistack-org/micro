package mock

import (
	"errors"
	"testing"

	"go.unistack.org/micro/v5/meter"
)

func TestMockMeter_Implements(t *testing.T) {
	var _ meter.Meter = NewMockMeter()
}

func TestMockMeter_Init_Success(t *testing.T) {
	m := NewMockMeter()
	m.ExpectInit()
	if err := m.Init(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMockMeter_Init_Error(t *testing.T) {
	m := NewMockMeter()
	m.ExpectInit().WillReturnError(errors.New("init failed"))
	err := m.Init()
	if err == nil || err.Error() != "init failed" {
		t.Fatalf("expected 'init failed', got %v", err)
	}
}

func TestMockMeter_Init_Unexpected(t *testing.T) {
	m := NewMockMeter()
	if err := m.Init(); err == nil {
		t.Fatal("expected error for unexpected Init call")
	}
}

func TestMockMeter_Counter(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectCounter("requests", "method", "GET")

	c := m.Counter("requests", "method", "GET")
	c.Inc()
	c.Inc()
	c.Add(3)
	c.Dec()
	c.Set(10)

	if got := exp.Counter().Get(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestMockMeter_Counter_Unexpected(t *testing.T) {
	m := NewMockMeter()
	_ = m.Counter("hits") // no expectation — goes to unexpected
	if len(m.unexpected) == 0 {
		t.Fatal("expected unexpected call to be recorded")
	}
}

func TestMockMeter_FloatCounter(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectFloatCounter("bytes")

	fc := m.FloatCounter("bytes")
	fc.Add(1.5)
	fc.Add(2.5)
	fc.Sub(1.0)
	fc.Set(5.0)

	if got := exp.FloatCounter().Get(); got != 5.0 {
		t.Fatalf("expected 5.0, got %f", got)
	}
}

func TestMockMeter_Gauge(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectGauge("temperature")

	g := m.Gauge("temperature", nil)
	g.Set(42.0)
	g.Inc()
	g.Dec()
	g.Add(0.5)

	if got := exp.Gauge().Get(); got != 42.5 {
		t.Fatalf("expected 42.5, got %f", got)
	}
}
