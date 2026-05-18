package ring

import (
	"testing"
	"time"
)

func TestBufferStream(t *testing.T) {
	b := New(10)

	entries, stop := b.Stream()
	if entries == nil {
		t.Fatal("expected entries channel, got nil")
	}
	if stop == nil {
		t.Fatal("expected stop channel, got nil")
	}

	b.Put("stream-entry")

	select {
	case e := <-entries:
		if e.Value.(string) != "stream-entry" {
			t.Fatalf("expected stream-entry got %v", e.Value)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for stream entry")
	}

	// stop the stream
	close(stop)
	// put another value to trigger cleanup
	b.Put("after-stop")
}

func TestBufferSize(t *testing.T) {
	b := New(5)
	if b.Size() != 5 {
		t.Fatalf("expected size 5 got %d", b.Size())
	}
}

func TestBufferSinceFuture(t *testing.T) {
	b := New(10)
	b.Put("entry1")

	// a time in the future should return nil
	future := time.Now().Add(time.Hour)
	v := b.Since(future)
	if v != nil {
		t.Fatalf("expected nil for future time, got %v", v)
	}
}

func TestBuffer(t *testing.T) {
	b := New(10)

	// test one value
	b.Put("foo")
	v := b.Get(1)

	if val := v[0].Value.(string); val != "foo" {
		t.Fatalf("expected foo got %v", val)
	}

	b = New(10)

	// test 10 values
	for i := range 10 {
		b.Put(i)
	}

	d := time.Now()
	v = b.Get(10)

	for i := range 10 {
		val := v[i].Value.(int)

		if val != i {
			t.Fatalf("expected %d got %d", i, val)
		}
	}

	// test more values

	for i := range 10 {
		v := i * 2
		b.Put(v)
	}

	v = b.Get(10)

	for i := range 10 {
		val := v[i].Value.(int)
		expect := i * 2
		if val != expect {
			t.Fatalf("expected %d got %d", expect, val)
		}
	}

	// sleep 100 ms
	time.Sleep(time.Millisecond * 100)

	// assume we'll get everything
	v = b.Since(d)

	if len(v) != 10 {
		t.Fatalf("expected 10 entries but got %d", len(v))
	}

	// write 1 more entry
	d = time.Now()
	b.Put(100)

	// sleep 100 ms
	time.Sleep(time.Millisecond * 100)

	v = b.Since(d)
	if len(v) != 1 {
		t.Fatalf("expected 1 entries but got %d", len(v))
	}

	if v[0].Value.(int) != 100 {
		t.Fatalf("expected value 100 got %v", v[0])
	}
}
