// Package codec is an interface for encoding messages
package codec

import (
	"errors"
	"io"

	"github.com/unistack-org/micro/v3/metadata"
)

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
	DefaultMaxMessageSize = 1024 * 1024 * 4 // 4Mb
)

// MessageType
type MessageType int

// Codec encodes/decodes various types of messages used within micro.
// ReadHeader and ReadBody are called in pairs to read requests/responses
// from the connection. Close is called when finished with the
// connection. ReadBody may be called with a nil argument to force the
// body to be read and discarded.
type Codec interface {
	ReadHeader(io.ReadWriter, *Message, MessageType) error
	ReadBody(io.ReadWriter, interface{}) error
	Write(io.ReadWriter, *Message, interface{}) error
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}

// Message represents detailed information about
// the communication, likely followed by the body.
// In the case of an error, body may be nil.
type Message struct {
	Id       string
	Type     MessageType
	Target   string
	Method   string
	Endpoint string
	Error    string

	// The values read from the socket
	Header metadata.Metadata
	Body   []byte
}

type Option func(*Options)

type Options struct {
	MaxMessageSize int64
}

func MaxMessageSize(n int64) Option {
	return func(o *Options) {
		o.MaxMessageSize = n
	}
}

func NewMessage(t MessageType) *Message {
	return &Message{Type: t, Header: metadata.New(0)}
}
