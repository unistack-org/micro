// Package codec is an interface for encoding messages
package codec

import (
	"errors"
	"io"

	"github.com/unistack-org/micro/v3/metadata"
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
	DefaultMaxMsgSize int = 1024 * 1024 * 4 // 4Mb
	// DefaultCodec is the global default codec
	DefaultCodec Codec = NewCodec()
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
	ReadHeader(io.Reader, *Message, MessageType) error
	ReadBody(io.Reader, interface{}) error
	Write(io.Writer, *Message, interface{}) error
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
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
	Id       string
	Body     []byte
	Type     MessageType
}

// NewMessage creates new codec message
func NewMessage(t MessageType) *Message {
	return &Message{Type: t, Header: metadata.New(0)}
}
