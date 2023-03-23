package options

import "testing"

func TestHooks_Append(t *testing.T) {
	fn1 := func() {}
	fn2 := func() {}
	hs := &Hooks{}
	hs.Append(fn1, fn2)
	if len(*hs) != 2 {
		t.Fatalf("unexpected Append error")
	}
}

func TestHooks_Replace(t *testing.T) {
	fn1 := func() {}
	fn2 := func() {}
	hs := &Hooks{}
	hs.Append(fn1, fn2, fn1)
	if len(*hs) != 3 {
		t.Fatalf("unexpected Append error")
	}
	hs.Replace(fn1, fn2)
	if len(*hs) != 2 {
		t.Fatalf("unexpected Replace error")
	}
}

func TestHooks_EachNext(t *testing.T) {
	n := 5
	fn1 := func() {
		n *= 2
	}
	fn2 := func() {
		n -= 10
	}
	hs := &Hooks{}
	hs.Append(fn1, fn2)

	hs.EachNext(func(h Hook) {
		h.(func())()
	})
	if n != 0 {
		t.Fatalf("unexpected EachNext")
	}
}

func TestHooks_EachPrev(t *testing.T) {
	n := 5
	fn1 := func() {
		n *= 2
	}
	fn2 := func() {
		n -= 10
	}
	hs := &Hooks{}
	hs.Append(fn2, fn1)

	hs.EachPrev(func(h Hook) {
		h.(func())()
	})
	if n != 0 {
		t.Fatalf("unexpected EachPrev")
	}
}
