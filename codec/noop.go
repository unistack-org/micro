package codec

import (
	"encoding/json"

	codecpb "go.unistack.org/micro-proto/v4/codec"
)

type noopCodec struct {
	opts Options
}

func (c *noopCodec) String() string {
	return "noop"
}

// NewCodec returns new noop codec
func NewCodec(opts ...Option) Codec {
	return &noopCodec{opts: NewOptions(opts...)}
}

func (c *noopCodec) Marshal(v interface{}, opts ...Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	switch ve := v.(type) {
	case string:
		return []byte(ve), nil
	case *string:
		return []byte(*ve), nil
	case *[]byte:
		return *ve, nil
	case []byte:
		return ve, nil
	case *Frame:
		return ve.Data, nil
	case *codecpb.Frame:
		return ve.Data, nil
	}

	return json.Marshal(v)
}

func (c *noopCodec) Unmarshal(d []byte, v interface{}, opts ...Option) error {
	if v == nil {
		return nil
	}
	switch ve := v.(type) {
	case *string:
		*ve = string(d)
		return nil
	case []byte:
		copy(ve, d)
		return nil
	case *[]byte:
		*ve = d
		return nil
	case *Frame:
		ve.Data = d
		return nil
	case *codecpb.Frame:
		ve.Data = d
		return nil
	}

	return json.Unmarshal(d, v)
}
