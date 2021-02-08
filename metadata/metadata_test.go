package metadata

import (
	"context"
	"testing"
)

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
		//fmt.Printf("k: %s, v: %s\n", k, v)
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
