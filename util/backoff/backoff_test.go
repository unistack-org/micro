package backoff

import (
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	t.Run("zero attempts", func(t *testing.T) {
		d := Do(0)
		if d < 0 {
			t.Fatalf("expected non-negative duration, got %v", d)
		}
	})
	t.Run("small attempts", func(t *testing.T) {
		d := Do(2)
		if d <= 0 {
			t.Fatalf("expected positive duration, got %v", d)
		}
	})
	t.Run("capped at 2 minutes when attempts=14", func(t *testing.T) {
		d := Do(14)
		if d != 2*time.Minute {
			t.Fatalf("expected 2m, got %v", d)
		}
	})
	t.Run("capped at 2 minutes when attempts=100", func(t *testing.T) {
		d := Do(100)
		if d != 2*time.Minute {
			t.Fatalf("expected 2m, got %v", d)
		}
	})
	t.Run("attempts=13 is within cap", func(t *testing.T) {
		d := Do(13)
		if d > 2*time.Minute {
			t.Fatalf("attempts=13 should be <=2m, got %v", d)
		}
	})
}
