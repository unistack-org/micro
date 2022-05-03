// Package codec is an interface for encoding messages
package codec // import "go.unistack.org/micro/v3/codec"

import (
	"errors"
	"io"

	"go.unistack.org/micro/v3/metadata"
)

// Message types
const (
	Error MessageType = iota
	Request
	Response
	Event
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

// MessageType specifies message type for codec
type MessageType int

// Codec encodes/decodes various types of messages used within micro.
// ReadHeader and ReadBody are called in pairs to read requests/responses
// from the connection. Close is called when finished with the
// connection. ReadBody may be called with a nil argument to force the
// body to be read and discarded.
type Codec interface {
	ReadHeader(r io.Reader, m *Message, mt MessageType) error
	ReadBody(r io.Reader, v interface{}) error
	Write(w io.Writer, m *Message, v interface{}) error
	Marshal(v interface{}, opts ...Option) ([]byte, error)
	Unmarshal(b []byte, v interface{}, opts ...Option) error
	String() string
}

// Message represents detailed information about
// the communication, likely followed by the body.
// In the case of an error, body may be nil.
type Message struct {
	Header   metadata.Metadata
	Target   string
	Method   string
	Endpoint string
	Error    string
	ID       string
	Body     []byte
	Type     MessageType
}

// NewMessage creates new codec message
func NewMessage(t MessageType) *Message {
	return &Message{Type: t, Header: metadata.New(0)}
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
