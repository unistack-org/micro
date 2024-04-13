package pool

import (
	"bytes"
	"strings"
	"testing"
)

func TestBytes(t *testing.T) {
	p := NewPool(func() *bytes.Buffer { return bytes.NewBuffer(nil) })
	b := p.Get()
	b.Write([]byte(`test`))
	if b.String() != "test" {
		t.Fatal("pool not works")
	}
	p.Put(b)
}

func TestStrings(t *testing.T) {
	p := NewPool(func() *strings.Builder { return &strings.Builder{} })
	b := p.Get()
	b.Write([]byte(`test`))
	if b.String() != "test" {
		t.Fatal("pool not works")
	}
	p.Put(b)
}
