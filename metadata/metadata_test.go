package metadata

import (
	"context"
	"testing"
)

func TestMetadataSetMultiple(t *testing.T) {
	md := New(4)
	md.Set("key1", "val1", "key2", "val2", "key3")

	if v, ok := md.Get("key1"); !ok || v != "val1" {
		t.Fatalf("invalid kv %#+v", md)
	}
	if v, ok := md.Get("key2"); !ok || v != "val2" {
		t.Fatalf("invalid kv %#+v", md)
	}
	if _, ok := md.Get("key3"); ok {
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
	if _, ok := md.Get("key1"); !ok {
		t.Fatal("key1 not found")
	}
}

func TestPairs(t *testing.T) {
	md, ok := Pairs("key1", "val1", "key2", "val2")
	if !ok {
		t.Fatal("odd number of kv")
	}
	if _, ok = md.Get("key1"); !ok {
		t.Fatal("key1 not found")
	}
}

func testCtx(ctx context.Context) {
	md := New(2)
	md.Set("Key1", "Val1_new")
	md.Set("Key3", "Val3")
	SetOutgoingContext(ctx, md)
}

func TestPassing(t *testing.T) {
	ctx := context.TODO()
	md1 := New(2)
	md1.Set("Key1", "Val1")
	md1.Set("Key2", "Val2")

	ctx = NewIncomingContext(ctx, md1)
	testCtx(ctx)
	md, ok := FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("missing metadata from outgoing context")
	}
	if v, ok := md.Get("Key1"); !ok || v != "Val1_new" {
		t.Fatalf("invalid metadata value %#+v", md)
	}
}

func TestMerge(t *testing.T) {
	omd := Metadata{
		"key1": "val1",
	}
	mmd := Metadata{
		"key2": "val2",
	}

	nmd := Merge(omd, mmd, true)
	if len(nmd) != 2 {
		t.Fatalf("merge failed: %v", nmd)
	}
}

func TestIterator(t *testing.T) {
	md := Metadata{
		"1Last":   "last",
		"2First":  "first",
		"3Second": "second",
	}

	iter := md.Iterator()
	var k, v string

	for iter.Next(&k, &v) {
		// fmt.Printf("k: %s, v: %s\n", k, v)
	}
}

func TestMedataCanonicalKey(t *testing.T) {
	md := New(1)
	md.Set("x-request-id", "12345")
	v, ok := md.Get("x-request-id")
	if !ok {
		t.Fatalf("failed to get x-request-id")
	} else if v != "12345" {
		t.Fatalf("invalid metadata value: %s != %s", "12345", v)
	}

	v, ok = md.Get("X-Request-Id")
	if !ok {
		t.Fatalf("failed to get x-request-id")
	} else if v != "12345" {
		t.Fatalf("invalid metadata value: %s != %s", "12345", v)
	}
	v, ok = md.Get("X-Request-ID")
	if !ok {
		t.Fatalf("failed to get x-request-id")
	} else if v != "12345" {
		t.Fatalf("invalid metadata value: %s != %s", "12345", v)
	}
}

func TestMetadataSet(t *testing.T) {
	md := New(1)

	md.Set("Key", "val")

	val, ok := md.Get("Key")
	if !ok {
		t.Fatal("key Key not found")
	}
	if val != "val" {
		t.Errorf("key Key with value val != %v", val)
	}
}

func TestMetadataDelete(t *testing.T) {
	md := Metadata{
		"Foo": "bar",
		"Baz": "empty",
	}

	md.Del("Baz")
	_, ok := md.Get("Baz")
	if ok {
		t.Fatal("key Baz not deleted")
	}
}

func TestNilContext(t *testing.T) {
	var ctx context.Context

	_, ok := FromContext(ctx)
	if ok {
		t.Fatal("nil context")
	}
}

func TestMetadataCopy(t *testing.T) {
	md := Metadata{
		"Foo": "bar",
		"Bar": "baz",
	}

	cp := Copy(md)

	for k, v := range md {
		if cv := cp[k]; cv != v {
			t.Fatalf("Got %s:%s for %s:%s", k, cv, k, v)
		}
	}
}

func TestMetadataContext(t *testing.T) {
	md := Metadata{
		"Foo": "bar",
	}

	ctx := NewContext(context.TODO(), md)

	emd, ok := FromContext(ctx)
	if !ok {
		t.Errorf("Unexpected error retrieving metadata, got %t", ok)
	}

	if emd["Foo"] != md["Foo"] {
		t.Errorf("Expected key: %s val: %s, got key: %s val: %s", "Foo", md["Foo"], "Foo", emd["Foo"])
	}

	if i := len(emd); i != 1 {
		t.Errorf("Expected metadata length 1 got %d", i)
	}
}
