// Package codec is an interface for encoding messages
package codec

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

type CodecV2 interface {
	Marshal(buf []byte, v interface{}, opts ...Option) ([]byte, error)
	Unmarshal(buf []byte, v interface{}, opts ...Option) error
	String() string
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

// MarshalYAML returns m as the JSON encoding of m.
func (m *RawMessage) MarshalYAML() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	} else if len(*m) == 0 {
		return []byte("null"), nil
	}
	return *m, nil
}

// UnmarshalYAML sets *m to a copy of data.
func (m *RawMessage) UnmarshalYAML(data []byte) error {
	if m == nil {
		return errors.New("RawMessage UnmarshalYAML on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}
