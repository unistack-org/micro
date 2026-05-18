package metadata

import (
	"context"
	"testing"
)

func TesSet(t *testing.T) {
	md := Pairs("key1", "val1", "key2", "val2")
	md.Set("key1", "val2", "val3")
	v := md.GetJoined("X-Request-Id")
	if v != "val2, val3" {
		t.Fatal("set not works")
	}
}

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
	ctx := context.WithValue(context.TODO(), metadataCurrentKeyVal, rawMetadata{md: New(0)})

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
	ctx := context.WithValue(context.TODO(), metadataIncomingKeyVal, rawMetadata{md: New(0)})

	c, ok := FromIncomingContext(ctx)
	if c == nil || !ok {
		t.Fatal("FromIncomingContext not works")
	}
}

func TestFromOutgoingContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), metadataOutgoingKeyVal, rawMetadata{md: New(0)})

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

func TestNewWithMetadata(t *testing.T) {
	m := map[string]string{"key1": "val1", "key2": "val2"}
	md := NewWithMetadata(m)
	if v := md.GetJoined("key1"); v != "val1" {
		t.Fatalf("NewWithMetadata: expected val1, got %s", v)
	}
	if v := md.GetJoined("key2"); v != "val2" {
		t.Fatalf("NewWithMetadata: expected val2, got %s", v)
	}
}

func TestJoin(t *testing.T) {
	md1 := Pairs("key1", "val1")
	md2 := Pairs("key2", "val2")
	joined := Join(md1, md2)
	if v := joined.GetJoined("key1"); v != "val1" {
		t.Fatalf("Join: expected val1, got %s", v)
	}
	if v := joined.GetJoined("key2"); v != "val2" {
		t.Fatalf("Join: expected val2, got %s", v)
	}
}

func TestMetadataCopyMethod(t *testing.T) {
	md := Metadata{"foo": []string{"bar"}}
	cp := md.Copy()
	if v := cp.GetJoined("foo"); v != "bar" {
		t.Fatalf("Copy method: expected bar, got %s", v)
	}
	// Ensure deep copy
	md["foo"][0] = "changed"
	if v := cp.GetJoined("foo"); v != "bar" {
		t.Fatal("Copy method: not a deep copy")
	}
}

func TestCopyTo(t *testing.T) {
	src := Metadata{"key1": []string{"val1"}}
	dst := New(2)
	src.CopyTo(dst)
	if v := dst.GetJoined("key1"); v != "val1" {
		t.Fatalf("CopyTo: expected val1, got %s", v)
	}
}

func TestLen(t *testing.T) {
	md := Pairs("key1", "val1", "key2", "val2")
	if md.Len() != 2 {
		t.Fatalf("Len: expected 2, got %d", md.Len())
	}
}

func TestAsMap(t *testing.T) {
	md := Pairs("key1", "val1", "key2", "val2")
	m := md.AsMap()
	if m["key1"] != "val1" {
		t.Fatalf("AsMap: expected val1, got %s", m["key1"])
	}
}

func TestAsHTTP1(t *testing.T) {
	md := Metadata{"x-request-id": []string{"abc"}}
	h := md.AsHTTP1()
	if len(h["X-Request-Id"]) == 0 || h["X-Request-Id"][0] != "abc" {
		t.Fatalf("AsHTTP1: expected X-Request-Id=abc, got %v", h)
	}
}

func TestAsHTTP2(t *testing.T) {
	md := Metadata{"X-Request-Id": []string{"abc"}}
	h := md.AsHTTP2()
	if len(h["x-request-id"]) == 0 || h["x-request-id"][0] != "abc" {
		t.Fatalf("AsHTTP2: expected x-request-id=abc, got %v", h)
	}
}

func TestAppend(t *testing.T) {
	md := New(1)
	md.Append("key1", "val1")
	md.Append("key1", "val2")
	vals := md.Get("key1")
	if len(vals) != 2 {
		t.Fatalf("Append: expected 2 values, got %d", len(vals))
	}
	// Append with no values should be a no-op
	md.Append("key1")
	if len(md.Get("key1")) != 2 {
		t.Fatal("Append: empty call should be no-op")
	}
}

func TestSetNoValues(t *testing.T) {
	md := New(1)
	md.Set("key1") // no-op
	if v := md.Get("key1"); v != nil {
		t.Fatal("Set with no values should be no-op")
	}
}

func TestPairsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Pairs with odd args should panic")
		}
	}()
	Pairs("key1") // nolint:staticcheck
}

func TestAppendContext(t *testing.T) {
	ctx := NewContext(context.TODO(), New(1))
	ctx = AppendContext(ctx, "key1", "val1")
	md, ok := FromContext(ctx)
	if !ok {
		t.Fatal("AppendContext: context not found")
	}
	if v := md.GetJoined("key1"); v != "val1" {
		t.Fatalf("AppendContext: expected val1, got %s", v)
	}
}

func TestAppendContextNoPrior(t *testing.T) {
	// AppendContext on a context with no existing metadata
	ctx := AppendContext(context.TODO(), "k", "v")
	md, ok := FromContext(ctx)
	if !ok {
		t.Fatal("AppendContext on empty context: not found")
	}
	if v := md.GetJoined("k"); v != "v" {
		t.Fatalf("AppendContext on empty context: expected v, got %s", v)
	}
}

func TestAppendContextPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("AppendContext with odd args should panic")
		}
	}()
	AppendContext(context.TODO(), "key1") // nolint:staticcheck
}

func TestAppendOutgoingContextPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("AppendOutgoingContext with odd args should panic")
		}
	}()
	AppendOutgoingContext(context.TODO(), "key1") // nolint:staticcheck
}

func TestAppendOutgoingContextWithExisting(t *testing.T) {
	ctx := NewOutgoingContext(context.TODO(), Pairs("key1", "val1"))
	ctx = AppendOutgoingContext(ctx, "key2", "val2")
	md, ok := FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("AppendOutgoingContext with existing: not found")
	}
	if v := md.GetJoined("key1"); v != "val1" {
		t.Fatalf("AppendOutgoingContext: expected key1=val1, got %s", v)
	}
	if v := md.GetJoined("key2"); v != "val2" {
		t.Fatalf("AppendOutgoingContext: expected key2=val2, got %s", v)
	}
}

func TestMustContext(t *testing.T) {
	ctx := NewContext(context.TODO(), Pairs("key1", "val1"))
	md := MustContext(ctx)
	if v := md.GetJoined("key1"); v != "val1" {
		t.Fatalf("MustContext: expected val1, got %s", v)
	}
}

func TestMustContextPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustContext should panic when metadata missing")
		}
	}()
	MustContext(context.TODO())
}

func TestMustIncomingContext(t *testing.T) {
	ctx := NewIncomingContext(context.TODO(), Pairs("key1", "val1"))
	md := MustIncomingContext(ctx)
	if v := md.GetJoined("key1"); v != "val1" {
		t.Fatalf("MustIncomingContext: expected val1, got %s", v)
	}
}

func TestMustIncomingContextPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustIncomingContext should panic when metadata missing")
		}
	}()
	MustIncomingContext(context.TODO())
}

func TestMustOutgoingContext(t *testing.T) {
	ctx := NewOutgoingContext(context.TODO(), Pairs("key1", "val1"))
	md := MustOutgoingContext(ctx)
	if v := md.GetJoined("key1"); v != "val1" {
		t.Fatalf("MustOutgoingContext: expected val1, got %s", v)
	}
}

func TestMustOutgoingContextPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustOutgoingContext should panic when metadata missing")
		}
	}()
	MustOutgoingContext(context.TODO())
}

func TestValueFromCurrentContext(t *testing.T) {
	md := Pairs("key1", "val1")
	ctx := NewContext(context.TODO(), md)
	v := ValueFromCurrentContext(ctx, "key1")
	if len(v) == 0 || v[0] != "val1" {
		t.Fatalf("ValueFromCurrentContext: expected val1, got %v", v)
	}
	// Case-insensitive lookup
	v = ValueFromCurrentContext(ctx, "KEY1")
	if len(v) == 0 || v[0] != "val1" {
		t.Fatalf("ValueFromCurrentContext case-insensitive: expected val1, got %v", v)
	}
	// Missing key
	v = ValueFromCurrentContext(ctx, "missing")
	if v != nil {
		t.Fatalf("ValueFromCurrentContext: expected nil for missing key, got %v", v)
	}
}

func TestValueFromCurrentContextNoMetadata(t *testing.T) {
	v := ValueFromCurrentContext(context.TODO(), "key1")
	if v != nil {
		t.Fatalf("ValueFromCurrentContext with no metadata: expected nil, got %v", v)
	}
}

func TestValueFromIncomingContext(t *testing.T) {
	md := Pairs("key1", "val1")
	ctx := NewIncomingContext(context.TODO(), md)
	v := ValueFromIncomingContext(ctx, "key1")
	if len(v) == 0 || v[0] != "val1" {
		t.Fatalf("ValueFromIncomingContext: expected val1, got %v", v)
	}
	// Case-insensitive lookup
	v = ValueFromIncomingContext(ctx, "KEY1")
	if len(v) == 0 || v[0] != "val1" {
		t.Fatalf("ValueFromIncomingContext case-insensitive: expected val1, got %v", v)
	}
	// Missing key
	v = ValueFromIncomingContext(ctx, "missing")
	if v != nil {
		t.Fatalf("ValueFromIncomingContext: expected nil for missing key, got %v", v)
	}
}

func TestValueFromIncomingContextNoMetadata(t *testing.T) {
	v := ValueFromIncomingContext(context.TODO(), "key1")
	if v != nil {
		t.Fatalf("ValueFromIncomingContext with no metadata: expected nil, got %v", v)
	}
}

func TestValueFromOutgoingContext(t *testing.T) {
	md := Pairs("key1", "val1")
	ctx := NewOutgoingContext(context.TODO(), md)
	v := ValueFromOutgoingContext(ctx, "key1")
	if len(v) == 0 || v[0] != "val1" {
		t.Fatalf("ValueFromOutgoingContext: expected val1, got %v", v)
	}
	// Case-insensitive lookup
	v = ValueFromOutgoingContext(ctx, "KEY1")
	if len(v) == 0 || v[0] != "val1" {
		t.Fatalf("ValueFromOutgoingContext case-insensitive: expected val1, got %v", v)
	}
	// Missing key
	v = ValueFromOutgoingContext(ctx, "missing")
	if v != nil {
		t.Fatalf("ValueFromOutgoingContext: expected nil for missing key, got %v", v)
	}
}

func TestValueFromOutgoingContextNoMetadata(t *testing.T) {
	v := ValueFromOutgoingContext(context.TODO(), "key1")
	if v != nil {
		t.Fatalf("ValueFromOutgoingContext with no metadata: expected nil, got %v", v)
	}
}

func TestGetJoinedMultiple(t *testing.T) {
	md := New(1)
	md["key1"] = []string{"a", "b", "c"}
	v := md.GetJoined("key1")
	if v != "a,b,c" {
		t.Fatalf("GetJoined: expected a,b,c, got %s", v)
	}
}

func TestDelMultipleFormats(t *testing.T) {
	md := Metadata{
		"x-request-id": []string{"1"},
		"X-Request-Id": []string{"2"},
		"X-Request-ID": []string{"3"},
	}
	md.Del("x-request-id")
	if md.Get("x-request-id") != nil {
		t.Fatal("Del: lowercase key not deleted")
	}
}

func TestFromContextWithAdded(t *testing.T) {
	ctx := NewContext(context.TODO(), Pairs("key1", "val1"))
	ctx = AppendContext(ctx, "key2", "val2")
	md, ok := FromContext(ctx)
	if !ok {
		t.Fatal("FromContext with added: not found")
	}
	if v := md.GetJoined("key1"); v != "val1" {
		t.Fatalf("FromContext with added: expected key1=val1, got %s", v)
	}
	if v := md.GetJoined("key2"); v != "val2" {
		t.Fatalf("FromContext with added: expected key2=val2, got %s", v)
	}
}

func TestFromIncomingContextWithAdded(t *testing.T) {
	// Use rawMetadata directly to test the added path
	ctx := context.WithValue(context.TODO(), metadataIncomingKeyVal, rawMetadata{
		md:    Pairs("key1", "val1"),
		added: [][]string{{"key2", "val2"}},
	})
	md, ok := FromIncomingContext(ctx)
	if !ok {
		t.Fatal("FromIncomingContext with added: not found")
	}
	if v := md.GetJoined("key2"); v != "val2" {
		t.Fatalf("FromIncomingContext with added: expected key2=val2, got %s", v)
	}
}
