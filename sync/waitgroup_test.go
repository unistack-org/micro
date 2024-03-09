package sync

import (
	"context"
	"testing"
	"time"
)

func TestWaitGroupContext(t *testing.T) {
	wg := NewWaitGroup()
	_ = t
	wg.Add(1)
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	wg.WaitContext(ctx)
}

func TestWaitGroupReuse(t *testing.T) {
	wg := NewWaitGroup()
	defer func() {
		if wg.Waiters() != 0 {
			t.Fatal("lost goroutines")
		}
	}()

	wg.Add(1)
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	wg.WaitContext(ctx)

	wg.Add(1)
	defer wg.Done()
	ctx, cancel = context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	wg.WaitContext(ctx)
}
