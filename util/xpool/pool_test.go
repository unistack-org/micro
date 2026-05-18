package pool

import (
	"bytes"
	"testing"
)

func TestNewPool(t *testing.T) {
	p := NewPool(func() string { return "hello" }, 128)
	v := p.Get()
	if v != "hello" {
		t.Fatalf("expected hello got %v", v)
	}
	p.Put(v)
}

func TestBytePoolCap(t *testing.T) {
	p := NewBytePool(512)
	if p.Cap() != 512 {
		t.Fatalf("expected cap 512 got %d", p.Cap())
	}
}

func TestBytePoolClose(t *testing.T) {
	p := NewBytePool(256)
	p.Close() // should not panic
}

func TestBytesPoolCap(t *testing.T) {
	p := NewBytesPool(512)
	if p.Cap() != 512 {
		t.Fatalf("expected cap 512 got %d", p.Cap())
	}
}

func TestBytesPoolStats(t *testing.T) {
	p := NewBytesPool(1024)
	b := p.Get()
	b.Write([]byte(`hello`))
	p.Put(b)
	st := p.Stats()
	if st.Put != 1 {
		t.Fatalf("expected put=1 got %d", st.Put)
	}
}

func TestBytesPoolClose(t *testing.T) {
	p := NewBytesPool(256)
	p.Close() // should not panic
}

func TestBytesPoolOversized(t *testing.T) {
	p := NewBytesPool(4)
	b := p.Get()
	// write more than the pool capacity
	for range 10 {
		b.Write([]byte(`abcdefghij`))
	}
	p.Put(b) // should be rejected (ret incremented)
	st := p.Stats()
	if st.Ret != 1 {
		t.Fatalf("expected ret=1 got %d", st.Ret)
	}
}

func TestStringsPoolCap(t *testing.T) {
	p := NewStringsPool(64)
	if p.Cap() != 64 {
		t.Fatalf("expected cap 64 got %d", p.Cap())
	}
}

func TestStringsPoolStats(t *testing.T) {
	p := NewStringsPool(64)
	b := p.Get()
	b.WriteString("test")
	p.Put(b)
	st := p.Stats()
	if st.Put != 1 {
		t.Fatalf("expected put=1 got %d", st.Put)
	}
}

func TestStringsPoolClose(t *testing.T) {
	p := NewStringsPool(64)
	p.Close() // should not panic
}

func TestStringsPoolOversized(t *testing.T) {
	p := NewStringsPool(4)
	b := p.Get()
	for range 10 {
		b.WriteString("abcdefghij")
	}
	p.Put(b) // oversized, should be rejected
	st := p.Stats()
	if st.Ret != 1 {
		t.Fatalf("expected ret=1 got %d", st.Ret)
	}
}

func TestByte(t *testing.T) {
	p := NewBytePool(1024)
	b := p.Get()
	copy(*b, []byte(`test`))
	if bytes.Equal(*b, []byte("test")) {
		t.Fatal("pool not works")
	}
	p.Put(b)
	b = p.Get()
	for range 1500 {
		*b = append(*b, []byte(`test`)...)
	}
	p.Put(b)
	st := p.Stats()
	if st.Get != 2 && st.Put != 2 && st.Mis != 1 && st.Ret != 1 {
		t.Fatalf("pool stats error %#+v", st)
	}
}

func TestBytes(t *testing.T) {
	p := NewBytesPool(1024)
	b := p.Get()
	b.Write([]byte(`test`))
	if b.String() != "test" {
		t.Fatal("pool not works")
	}
	p.Put(b)
}

func TestStrings(t *testing.T) {
	p := NewStringsPool(20)
	b := p.Get()
	b.Write([]byte(`test`))
	if b.String() != "test" {
		t.Fatal("pool not works")
	}
	p.Put(b)
}
