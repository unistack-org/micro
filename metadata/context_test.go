package metadata

import (
	"context"
	"testing"
)

func TestFromNilContext(t *testing.T) {
	// nolint: staticcheck
	c, ok := FromContext(nil)
	if ok || c != nil {
		t.Fatal("FromContext not works")
	}
}

func TestNewNilContext(t *testing.T) {
	// nolint: staticcheck
	ctx := NewContext(nil, New(0))

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), mdKey{}, &rawMetadata{New(0)})

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromContext not works")
	}
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), New(0))

	c, ok := FromContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewContext not works")
	}
}

func TestFromIncomingContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), mdIncomingKey{}, &rawMetadata{New(0)})

	c, ok := FromIncomingContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromIncomingContext not works")
	}
}

func TestFromOutgoingContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), mdOutgoingKey{}, &rawMetadata{New(0)})

	c, ok := FromOutgoingContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromOutgoingContext not works")
	}
}

func TestSetIncomingContext(t *testing.T) {
	md := New(1)
	md.Set("key", "val")
	ctx := context.WithValue(context.TODO(), mdIncomingKey{}, &rawMetadata{})
	if !SetIncomingContext(ctx, md) {
		t.Fatal("SetIncomingContext not works")
	}
	md, ok := FromIncomingContext(ctx)
	if md == nil || !ok {
		t.Fatal("SetIncomingContext not works")
	} else if v, ok := md.Get("key"); !ok || v != "val" {
		t.Fatal("SetIncomingContext not works")
	}
}

func TestSetOutgoingContext(t *testing.T) {
	md := New(1)
	md.Set("key", "val")
	ctx := context.WithValue(context.TODO(), mdOutgoingKey{}, &rawMetadata{})
	if !SetOutgoingContext(ctx, md) {
		t.Fatal("SetOutgoingContext not works")
	}
	md, ok := FromOutgoingContext(ctx)
	if md == nil || !ok {
		t.Fatal("SetOutgoingContext not works")
	} else if v, ok := md.Get("key"); !ok || v != "val" {
		t.Fatal("SetOutgoingContext not works")
	}
}

func TestNewIncomingContext(t *testing.T) {
	md := New(1)
	md.Set("key", "val")
	ctx := NewIncomingContext(context.TODO(), md)

	c, ok := FromIncomingContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewIncomingContext not works")
	}
}

func TestNewOutgoingContext(t *testing.T) {
	md := New(1)
	md.Set("key", "val")
	ctx := NewOutgoingContext(context.TODO(), md)

	c, ok := FromOutgoingContext(ctx)
	if c == nil || !ok {
		t.Fatal("NewOutgoingContext not works")
	}
}

func TestAppendIncomingContext(t *testing.T) {
	md := New(1)
	md.Set("key1", "val1")
	ctx := AppendIncomingContext(context.TODO(), "key2", "val2")

	nmd, ok := FromIncomingContext(ctx)
	if nmd == nil || !ok {
		t.Fatal("AppendIncomingContext not works")
	}
	if v, ok := nmd.Get("key2"); !ok || v != "val2" {
		t.Fatal("AppendIncomingContext not works")
	}
}

func TestAppendOutgoingContext(t *testing.T) {
	md := New(1)
	md.Set("key1", "val1")
	ctx := AppendOutgoingContext(context.TODO(), "key2", "val2")

	nmd, ok := FromOutgoingContext(ctx)
	if nmd == nil || !ok {
		t.Fatal("AppendOutgoingContext not works")
	}
	if v, ok := nmd.Get("key2"); !ok || v != "val2" {
		t.Fatal("AppendOutgoingContext not works")
	}
}
