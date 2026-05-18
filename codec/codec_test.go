package codec

import (
	"context"
	"testing"

	codecpb "go.unistack.org/micro-proto/v5/codec"
)

// TestRawMessageMarshalJSON tests RawMessage.MarshalJSON
func TestRawMessageMarshalJSON(t *testing.T) {
	// nil pointer
	var nilMsg *RawMessage
	b, err := nilMsg.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "null" {
		t.Fatalf("expected null, got %s", b)
	}

	// empty slice
	emptyMsg := RawMessage{}
	b, err = emptyMsg.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "null" {
		t.Fatalf("expected null, got %s", b)
	}

	// non-empty
	msg := RawMessage(`{"key":"value"}`)
	b, err = msg.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"key":"value"}` {
		t.Fatalf("unexpected value: %s", b)
	}
}

// TestRawMessageUnmarshalJSON tests RawMessage.UnmarshalJSON
func TestRawMessageUnmarshalJSON(t *testing.T) {
	// nil pointer error
	var nilMsg *RawMessage
	if err := nilMsg.UnmarshalJSON([]byte(`{}`)); err == nil {
		t.Fatal("expected error for nil RawMessage")
	}

	// normal case
	var msg RawMessage
	data := []byte(`{"key":"value"}`)
	if err := msg.UnmarshalJSON(data); err != nil {
		t.Fatal(err)
	}
	if string(msg) != string(data) {
		t.Fatalf("expected %s, got %s", data, msg)
	}
}

// TestRawMessageMarshalYAML tests RawMessage.MarshalYAML
func TestRawMessageMarshalYAML(t *testing.T) {
	// nil pointer
	var nilMsg *RawMessage
	b, err := nilMsg.MarshalYAML()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "null" {
		t.Fatalf("expected null, got %s", b)
	}

	// empty slice
	emptyMsg := RawMessage{}
	b, err = emptyMsg.MarshalYAML()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "null" {
		t.Fatalf("expected null, got %s", b)
	}

	// non-empty
	msg := RawMessage(`key: value`)
	b, err = msg.MarshalYAML()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `key: value` {
		t.Fatalf("unexpected value: %s", b)
	}
}

// TestRawMessageUnmarshalYAML tests RawMessage.UnmarshalYAML
func TestRawMessageUnmarshalYAML(t *testing.T) {
	// nil pointer error
	var nilMsg *RawMessage
	if err := nilMsg.UnmarshalYAML([]byte(`key: value`)); err == nil {
		t.Fatal("expected error for nil RawMessage")
	}

	// normal case
	var msg RawMessage
	data := []byte(`key: value`)
	if err := msg.UnmarshalYAML(data); err != nil {
		t.Fatal(err)
	}
	if string(msg) != string(data) {
		t.Fatalf("expected %s, got %s", data, msg)
	}
}

// TestNoopCodecString tests String() method
func TestNoopCodecString(t *testing.T) {
	nc := NewCodec()
	if nc.String() != "noop" {
		t.Fatalf("expected noop, got %s", nc.String())
	}
}

// TestNoopCodecMarshalNil tests Marshal with nil value
func TestNoopCodecMarshalNil(t *testing.T) {
	nc := NewCodec()
	b, err := nc.Marshal(nil)
	if err != nil {
		t.Fatal(err)
	}
	if b != nil {
		t.Fatalf("expected nil, got %v", b)
	}
}

// TestNoopCodecMarshalString tests Marshal with string
func TestNoopCodecMarshalString(t *testing.T) {
	nc := NewCodec()
	b, err := nc.Marshal("hello")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Fatalf("expected hello, got %s", b)
	}
}

// TestNoopCodecMarshalStringPtr tests Marshal with *string
func TestNoopCodecMarshalStringPtr(t *testing.T) {
	nc := NewCodec()
	s := "world"
	b, err := nc.Marshal(&s)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "world" {
		t.Fatalf("expected world, got %s", b)
	}
}

// TestNoopCodecMarshalBytesPtr tests Marshal with *[]byte
func TestNoopCodecMarshalBytesPtr(t *testing.T) {
	nc := NewCodec()
	data := []byte("bytes")
	b, err := nc.Marshal(&data)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "bytes" {
		t.Fatalf("expected bytes, got %s", b)
	}
}

// TestNoopCodecMarshalBytes tests Marshal with []byte
func TestNoopCodecMarshalBytes(t *testing.T) {
	nc := NewCodec()
	b, err := nc.Marshal([]byte("raw"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "raw" {
		t.Fatalf("expected raw, got %s", b)
	}
}

// TestNoopCodecMarshalFrame tests Marshal with *codecpb.Frame
func TestNoopCodecMarshalFrame(t *testing.T) {
	nc := NewCodec()
	frame := &codecpb.Frame{Data: []byte("frame data")}
	b, err := nc.Marshal(frame)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "frame data" {
		t.Fatalf("expected frame data, got %s", b)
	}
}

// TestNoopCodecMarshalJSON tests Marshal falling back to JSON
func TestNoopCodecMarshalJSON(t *testing.T) {
	nc := NewCodec()
	type payload struct {
		Name string `json:"name"`
	}
	b, err := nc.Marshal(&payload{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"name":"test"}` {
		t.Fatalf("unexpected json: %s", b)
	}
}

// TestNoopCodecUnmarshalNil tests Unmarshal with nil value
func TestNoopCodecUnmarshalNil(t *testing.T) {
	nc := NewCodec()
	if err := nc.Unmarshal([]byte("data"), nil); err != nil {
		t.Fatal(err)
	}
}

// TestNoopCodecUnmarshalBytesSlice tests Unmarshal into plain []byte (copy path)
func TestNoopCodecUnmarshalBytesSlice(t *testing.T) {
	nc := NewCodec()
	src := []byte("hello")
	dst := make([]byte, len(src))
	if err := nc.Unmarshal(src, dst); err != nil {
		t.Fatal(err)
	}
	if string(dst) != "hello" {
		t.Fatalf("expected hello, got %s", dst)
	}
}

// TestNoopCodecUnmarshalFrame tests Unmarshal into *codecpb.Frame
func TestNoopCodecUnmarshalFrame(t *testing.T) {
	nc := NewCodec()
	frame := &codecpb.Frame{}
	if err := nc.Unmarshal([]byte("frame"), frame); err != nil {
		t.Fatal(err)
	}
	if string(frame.Data) != "frame" {
		t.Fatalf("expected frame, got %s", frame.Data)
	}
}

// TestNoopCodecUnmarshalJSON tests Unmarshal falling back to JSON
func TestNoopCodecUnmarshalJSON(t *testing.T) {
	nc := NewCodec()
	type payload struct {
		Name string `json:"name"`
	}
	var p payload
	if err := nc.Unmarshal([]byte(`{"name":"json"}`), &p); err != nil {
		t.Fatal(err)
	}
	if p.Name != "json" {
		t.Fatalf("expected json, got %s", p.Name)
	}
}

// TestNewOptions tests NewOptions with various option setters
func TestNewOptions(t *testing.T) {
	opts := NewOptions(
		TagName("mytag"),
		Flatten(true),
	)
	if opts.TagName != "mytag" {
		t.Fatalf("expected mytag, got %s", opts.TagName)
	}
	if !opts.Flatten {
		t.Fatal("expected Flatten=true")
	}
}

// TestOptionLogger tests Logger option
func TestOptionLogger(t *testing.T) {
	opts := NewOptions(Logger(nil))
	if opts.Logger != nil {
		t.Fatal("expected nil logger")
	}
}

// TestOptionTracer tests Tracer option
func TestOptionTracer(t *testing.T) {
	opts := NewOptions(Tracer(nil))
	if opts.Tracer != nil {
		t.Fatal("expected nil tracer")
	}
}

// TestOptionMeter tests Meter option
func TestOptionMeter(t *testing.T) {
	opts := NewOptions(Meter(nil))
	if opts.Meter != nil {
		t.Fatal("expected nil meter")
	}
}

// TestMustContext tests MustContext panics when no codec
func TestMustContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from MustContext")
		}
	}()
	MustContext(context.Background())
}

// TestMustContextOK tests MustContext returns codec when present
func TestMustContextOK(t *testing.T) {
	nc := NewCodec()
	ctx := NewContext(context.Background(), nc)
	c := MustContext(ctx)
	if c == nil {
		t.Fatal("expected codec from MustContext")
	}
}

// TestFromContextNil tests FromContext with nil context
func TestFromContextNil(t *testing.T) {
	c, ok := FromContext(nil) // nolint: staticcheck
	if ok || c != nil {
		t.Fatal("expected nil codec and false for nil context")
	}
}

// TestNewContextNil tests NewContext with nil context
func TestNewContextNil(t *testing.T) {
	nc := NewCodec()
	ctx := NewContext(nil, nc) // nolint:staticcheck
	c, ok := FromContext(ctx)
	if !ok || c == nil {
		t.Fatal("expected codec from NewContext(nil, ...)")
	}
}

// TestDefaultCodec tests DefaultCodec variable
func TestDefaultCodec(t *testing.T) {
	if DefaultCodec == nil {
		t.Fatal("DefaultCodec should not be nil")
	}
	if DefaultCodec.String() != "noop" {
		t.Fatalf("expected noop, got %s", DefaultCodec.String())
	}
}

// TestDefaultTagName tests DefaultTagName variable
func TestDefaultTagName(t *testing.T) {
	if DefaultTagName != "codec" {
		t.Fatalf("expected codec, got %s", DefaultTagName)
	}
}
