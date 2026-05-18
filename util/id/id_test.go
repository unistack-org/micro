package id

import "testing"

func TestUUIDv8(t *testing.T) {
	id, err := New()
	if err != nil {
		t.Fatal(err)
	}
	_ = id
}

func TestToUUID(t *testing.T) {
	id, err := New()
	if err != nil {
		t.Fatal(err)
	}
	u := ToUUID(id)
	_ = u
}

func TestTypeNanoid(t *testing.T) {
	id, err := New(
		func(o *Options) { o.Type = TypeNanoid },
		WithNanoidAlphabet(DefaultNanoidAlphabet),
		WithNanoidSize(DefaultNanoidSize),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(id) == 0 {
		t.Fatal("expected non-empty id")
	}
}

func TestTypeNanoidEmptyAlphabet(t *testing.T) {
	_, err := New(
		func(o *Options) { o.Type = TypeNanoid },
		WithNanoidAlphabet(""),
	)
	if err == nil {
		t.Fatal("expected error for empty alphabet")
	}
}

func TestTypeNanoidNegativeSize(t *testing.T) {
	_, err := New(
		func(o *Options) { o.Type = TypeNanoid },
		WithNanoidAlphabet(DefaultNanoidAlphabet),
		WithNanoidSize(-1),
	)
	if err == nil {
		t.Fatal("expected error for negative size")
	}
}

func TestTypeUUIDv7(t *testing.T) {
	id, err := New(func(o *Options) { o.Type = TypeUUIDv7 })
	if err != nil {
		t.Fatal(err)
	}
	if len(id) == 0 {
		t.Fatal("expected non-empty id")
	}
}

func TestTypeUUIDv8Error(t *testing.T) {
	_, err := New(func(o *Options) { o.Type = TypeUUIDv8 })
	if err == nil {
		t.Fatal("expected error for UUIDv8")
	}
}

func TestTypeUnspecified(t *testing.T) {
	_, err := New(func(o *Options) { o.Type = TypeUnspecified })
	if err == nil {
		t.Fatal("expected error for unspecified type")
	}
}

func TestMustNew(t *testing.T) {
	id := MustNew()
	if len(id) == 0 {
		t.Fatal("expected non-empty id")
	}
}

func TestWithUUIDNode(t *testing.T) {
	var node [6]byte
	copy(node[:], []byte{1, 2, 3, 4, 5, 6})
	opts := NewOptions(WithUUIDNode(node))
	if opts.UUIDNode != node {
		t.Fatal("UUIDNode not set correctly")
	}
}

func TestGeneratorNew(t *testing.T) {
	g := &Generator{opts: NewOptions()}
	id, err := g.New()
	if err != nil {
		t.Fatal(err)
	}
	if len(id) == 0 {
		t.Fatal("expected non-empty id")
	}
}

func TestGeneratorMustNew(t *testing.T) {
	g := &Generator{opts: NewOptions()}
	id := g.MustNew()
	if len(id) == 0 {
		t.Fatal("expected non-empty id")
	}
}

func TestGeneratorNanoid(t *testing.T) {
	g := &Generator{opts: NewOptions(
		func(o *Options) { o.Type = TypeNanoid },
		WithNanoidAlphabet(DefaultNanoidAlphabet),
		WithNanoidSize(DefaultNanoidSize),
	)}
	id, err := g.New()
	if err != nil {
		t.Fatal(err)
	}
	if len(id) == 0 {
		t.Fatal("expected non-empty id")
	}
}

func TestGeneratorUUIDv8(t *testing.T) {
	g := &Generator{opts: NewOptions(func(o *Options) { o.Type = TypeUUIDv8 })}
	_, err := g.New()
	if err == nil {
		t.Fatal("expected error for UUIDv8")
	}
}

func TestGeneratorUnspecified(t *testing.T) {
	g := &Generator{opts: NewOptions(func(o *Options) { o.Type = TypeUnspecified })}
	_, err := g.New()
	if err == nil {
		t.Fatal("expected error for unspecified type")
	}
}
