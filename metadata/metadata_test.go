package metadata

import (
	"context"
	"testing"
)

/*
func TestAppendOutgoingContextModify(t *testing.T) {
	md := Pairs("key1", "val1")
	ctx := NewOutgoingContext(context.TODO(), md)
	nctx := AppendOutgoingContext(ctx, "key1", "val3", "key2", "val2")
	_ = nctx
	omd := MustOutgoingContext(nctx)
	fmt.Printf("%#+v\n", omd)
}
*/

func TestLowercase(t *testing.T) {
	md := New(1)
	md["x-request-id"] = []string{"12345"}
	v := md.GetJoined("X-Request-Id")
	if v == "" {
		t.Fatalf("metadata invalid %#+v", md)
	}
}

func TestMultipleUsage(t *testing.T) {
	ctx := context.TODO()
	md := New(0)
	md.Set("key1_1", "val1_1", "key1_2", "val1_2", "key1_3", "val1_3")
	ctx = NewIncomingContext(ctx, Copy(md))
	ctx = NewOutgoingContext(ctx, Copy(md))
	imd, _ := FromIncomingContext(ctx)
	omd, _ := FromOutgoingContext(ctx)
	_ = func(x context.Context) context.Context {
		m, _ := FromIncomingContext(x)
		m.Del("key1_2")
		return ctx
	}(ctx)
	_ = func(x context.Context) context.Context {
		m, _ := FromIncomingContext(x)
		m.Del("key1_3")
		return ctx
	}(ctx)
	_ = imd
	_ = omd
}

func TestMetadataSetMultiple(t *testing.T) {
	md := New(4)
	md.Set("key1", "val1", "key2", "val2")

	if v := md.GetJoined("key1"); v != "val1" {
		t.Fatalf("invalid kv %#+v", md)
	}
	if v := md.GetJoined("key2"); v != "val2" {
		t.Fatalf("invalid kv %#+v", md)
	}
}

func TestAppend(t *testing.T) {
	ctx := context.Background()
	ctx = AppendIncomingContext(ctx, "key1", "val1", "key2", "val2")
	md, ok := FromIncomingContext(ctx)
	if !ok {
		t.Fatal("metadata empty")
	}
	if v := md.Get("key1"); v == nil {
		t.Fatal("key1 not found")
	}
}

func TestPairs(t *testing.T) {
	md := Pairs("key1", "val1", "key2", "val2")
	if v := md.Get("key1"); v == nil {
		t.Fatal("key1 not found")
	}
}

func TestPassing(t *testing.T) {
	ctx := context.TODO()
	md1 := New(2)
	md1.Set("Key1", "Val1")
	md1.Set("Key2", "Val2")

	ctx = NewIncomingContext(ctx, md1)

	_, ok := FromOutgoingContext(ctx)
	if ok {
		t.Fatalf("create outgoing context")
	}

	ctx = NewOutgoingContext(ctx, md1)

	md, ok := FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("missing metadata from outgoing context")
	}
	if v := md.Get("Key1"); v == nil || v[0] != "Val1" {
		t.Fatalf("invalid metadata value %#+v", md)
	}
}

func TestIterator(t *testing.T) {
	md := Pairs(
		"1Last", "last",
		"2First", "first",
		"3Second", "second",
	)

	iter := md.Iterator()
	var k string
	var v []string
	chk := New(3)
	for iter.Next(&k, &v) {
		chk[k] = v
	}

	for k, v := range chk {
		if cv, ok := md[k]; !ok || len(cv) != len(v) || cv[0] != v[0] {
			t.Fatalf("XXXX %#+v %#+v", chk, md)
		}
	}
}

func TestMedataCanonicalKey(t *testing.T) {
	md := New(1)
	md.Set("x-request-id", "12345")
	v := md.GetJoined("x-request-id")
	if v == "" {
		t.Fatalf("failed to get x-request-id")
	} else if v != "12345" {
		t.Fatalf("invalid metadata value: %s != %s", "12345", v)
	}

	v = md.GetJoined("X-Request-Id")
	if v == "" {
		t.Fatalf("failed to get x-request-id")
	} else if v != "12345" {
		t.Fatalf("invalid metadata value: %s != %s", "12345", v)
	}
	v = md.GetJoined("X-Request-ID")
	if v == "" {
		t.Fatalf("failed to get x-request-id")
	} else if v != "12345" {
		t.Fatalf("invalid metadata value: %s != %s", "12345", v)
	}
}

func TestMetadataSet(t *testing.T) {
	md := New(1)

	md.Set("Key", "val")

	val := md.GetJoined("Key")
	if val == "" {
		t.Fatal("key Key not found")
	}
	if val != "val" {
		t.Errorf("key Key with value val != %v", val)
	}
}

func TestMetadataDelete(t *testing.T) {
	md := Metadata{
		"Foo": []string{"bar"},
		"Baz": []string{"empty"},
	}

	md.Del("Baz")
	v := md.Get("Baz")
	if v != nil {
		t.Fatal("key Baz not deleted")
	}
}

func TestMetadataCopy(t *testing.T) {
	md := Metadata{
		"Foo": []string{"bar"},
		"Bar": []string{"baz"},
	}

	cp := Copy(md)

	for k, v := range md {
		if cv := cp[k]; cv[0] != v[0] {
			t.Fatalf("Got %s:%s for %s:%s", k, cv, k, v)
		}
	}
}

func TestMetadataContext(t *testing.T) {
	md := Metadata{
		"Foo": []string{"bar"},
	}

	ctx := NewContext(context.TODO(), md)

	emd, ok := FromContext(ctx)
	if !ok {
		t.Errorf("Unexpected error retrieving metadata, got %t", ok)
	}

	if emd["Foo"][0] != md["Foo"][0] {
		t.Errorf("Expected key: %s val: %s, got key: %s val: %s", "Foo", md["Foo"], "Foo", emd["Foo"])
	}

	if i := len(emd); i != 1 {
		t.Errorf("Expected metadata length 1 got %d", i)
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), metadataCurrentKey{}, rawMetadata{md: New(0)})

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
	ctx := context.WithValue(context.TODO(), metadataIncomingKey{}, rawMetadata{md: New(0)})

	c, ok := FromIncomingContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromIncomingContext not works")
	}
}

func TestFromOutgoingContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), metadataOutgoingKey{}, rawMetadata{md: New(0)})

	c, ok := FromOutgoingContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromOutgoingContext not works")
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
	if v := nmd.GetJoined("key2"); v != "val2" {
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
	if v := nmd.GetJoined("key2"); v != "val2" {
		t.Fatal("AppendOutgoingContext not works")
	}
}
