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
