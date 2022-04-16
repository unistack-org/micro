package client

import (
	"context"
	"testing"
	"time"
)

func TestBackoffExp(t *testing.T) {
	results := []time.Duration{
		0 * time.Second,
		100 * time.Millisecond,
		600 * time.Millisecond,
		1900 * time.Millisecond,
		4300 * time.Millisecond,
		7900 * time.Millisecond,
	}

	r := &testRequest{
		service: "test",
		method:  "test",
	}

	for i := 0; i < 5; i++ {
		d, err := BackoffExp(context.TODO(), r, i)
		if err != nil {
			t.Fatal(err)
		}

		if d != results[i] {
			t.Fatalf("Expected equal than %v, got %v", results[i], d)
		}
	}
}

func TestBackoffInterval(t *testing.T) {
	min := 100 * time.Millisecond
	max := 300 * time.Millisecond

	r := &testRequest{
		service: "test",
		method:  "test",
	}

	fn := BackoffInterval(min, max)
	for i := 0; i < 5; i++ {
		d, err := fn(context.TODO(), r, i)
		if err != nil {
			t.Fatal(err)
		}

		if d < min || d > max {
			t.Fatalf("Expected %v < %v < %v", min, d, max)
		}
	}
}
