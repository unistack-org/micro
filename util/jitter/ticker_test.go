package jitter

import (
	"context"
	"testing"
	"time"
)

func TestNewTickerContext(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	ticker := NewTickerContext(ctx, 600*time.Millisecond, 1000*time.Millisecond)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			break loop
		case v, ok := <-ticker.C:
			if ok {
				t.Fatalf("context must be closed %s", v)
			}
			break loop
		}
	}
}

func TestTicker(t *testing.T) {
	t.Parallel()

	min := time.Duration(10)
	max := time.Duration(20)

	// tick can take a little longer since we're not adjusting it to account for
	// processing.
	precision := time.Duration(4)

	rt := NewTicker(min*time.Millisecond, max*time.Millisecond)
	for i := 0; i < 5; i++ {
		t0 := time.Now()
		t1 := <-rt.C
		td := t1.Sub(t0)
		if td < min*time.Millisecond {
			t.Fatalf("tick was shorter than expected: %s", td)
		} else if td > (max+precision)*time.Millisecond {
			t.Fatalf("tick was longer than expected: %s", td)
		}
	}
	rt.Stop()
	time.Sleep((max + precision) * time.Millisecond)
	select {
	case v, ok := <-rt.C:
		if ok || !v.IsZero() {
			t.Fatal("ticker did not shut down")
		}
	default:
		t.Fatal("expected to receive close channel signal")
	}
}
