// Package codec is an interface for encoding messages
package codec // import "go.unistack.org/micro/v3/codec"

import (
	"errors"
)

var (
	// ErrInvalidMessage returned when invalid messge passed to codec
	ErrInvalidMessage = errors.New("invalid message")
	// ErrUnknownContentType returned when content-type is unknown
	ErrUnknownContentType = errors.New("unknown content-type")
)

var (
	// DefaultMaxMsgSize specifies how much data codec can handle
	DefaultMaxMsgSize = 1024 * 1024 * 4 // 4Mb
	// DefaultCodec is the global default codec
	DefaultCodec = NewCodec()
	// DefaultTagName specifies struct tag name to control codec Marshal/Unmarshal
	DefaultTagName = "codec"
)

// Codec encodes/decodes various types of messages.
type Codec interface {
	Marshal(v interface{}, opts ...Option) ([]byte, error)
	Unmarshal(b []byte, v interface{}, opts ...Option) error
	String() string
}

// MarshalAppend calls codec.Marshal(v) and returns the data appended to buf.
// If codec implements MarshalAppend, that is called instead.
func MarshalAppend(buf []byte, c Codec, v interface{}, opts ...Option) ([]byte, error) {
	if nc, ok := c.(interface {
		MarshalAppend([]byte, interface{}, ...Option) ([]byte, error)
	}); ok {
		return nc.MarshalAppend(buf, v, opts...)
	}

	mbuf, err := c.Marshal(v, opts...)
	if err != nil {
		return nil, err
	}

	return append(buf, mbuf...), nil
}

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can be used to delay decoding or precompute a encoding.
type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m *RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	} else if len(*m) == 0 {
		return []byte("null"), nil
	}
	return *m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("RawMessage UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}
