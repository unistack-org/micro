package json

import (
	"bytes"
	"encoding/json"

	oldjsonpb "github.com/golang/protobuf/jsonpb"
	oldproto "github.com/golang/protobuf/proto"
	"github.com/oxtoacart/bpool"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var jsonpbMarshaler = &jsonpb.MarshalOptions{}
var oldjsonpbMarshaler = &oldjsonpb.Marshaler{}

// create buffer pool with 16 instances each preallocated with 256 bytes
var bufferPool = bpool.NewSizedBufferPool(16, 256)

type Marshaler struct{}

func (j Marshaler) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case proto.Message:
		return jsonpbMarshaler.Marshal(m)
	case oldproto.Message:
		buf, err := oldjsonpbMarshaler.MarshalToString(m)
		return []byte(buf), err
	}
	return json.Marshal(v)
}

func (j Marshaler) Unmarshal(d []byte, v interface{}) error {
	switch m := v.(type) {
	case proto.Message:
		return jsonpb.Unmarshal(d, m)
	case oldproto.Message:
		return oldjsonpb.Unmarshal(bytes.NewReader(d), m)
	}
	return json.Unmarshal(d, v)
}

func (j Marshaler) String() string {
	return "json"
}
