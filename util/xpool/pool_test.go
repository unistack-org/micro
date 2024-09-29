package pool

import (
	"bytes"
	"testing"
)

func TestByte(t *testing.T) {
	p := NewBytePool(1024)
	b := p.Get()
	copy(*b, []byte(`test`))
	if bytes.Equal(*b, []byte("test")) {
		t.Fatal("pool not works")
	}
	p.Put(b)
	b = p.Get()
	for i := 0; i < 1500; i++ {
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
