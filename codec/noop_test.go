package codec

import (
	"bytes"
	"testing"
)

func TestNoopBytes(t *testing.T) {
	req := []byte("test req")
	rsp := make([]byte, len(req))

	nc := NewCodec()
	if err := nc.Unmarshal(req, &rsp); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(req, rsp) {
		t.Fatalf("req not eq rsp: %s != %s", req, rsp)
	}
}

func TestNoopString(t *testing.T) {
	req := []byte("test req")
	var rsp string

	nc := NewCodec()
	if err := nc.Unmarshal(req, &rsp); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(req, []byte(rsp)) {
		t.Fatalf("req not eq rsp: %s != %s", req, rsp)
	}
}
