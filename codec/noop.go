package codec

import (
	"io"
	"io/ioutil"
)

type noopCodec struct {
}

// Frame gives us the ability to define raw data to send over the pipes
type Frame struct {
	Data []byte
}

func (c *noopCodec) ReadHeader(conn io.Reader, m *Message, t MessageType) error {
	return nil
}

func (c *noopCodec) ReadBody(conn io.Reader, b interface{}) error {
	// read bytes
	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		return err
	}

	if b == nil {
		return nil
	}

	switch v := b.(type) {
	case string:
		v = string(buf)
	case *string:
		*v = string(buf)
	case []byte:
		v = buf
	case *[]byte:
		*v = buf
	case *Frame:
		v.Data = buf
	default:
		return ErrInvalidMessage
	}

	return nil
}

func (c *noopCodec) Write(conn io.Writer, m *Message, b interface{}) error {
	if b == nil {
		return nil
	}

	var v []byte
	switch vb := b.(type) {
	case *Frame:
		v = vb.Data
	case string:
		v = []byte(vb)
	case *string:
		v = []byte(*vb)
	case *[]byte:
		v = *vb
	case []byte:
		v = vb
	default:
		return ErrInvalidMessage
	}
	_, err := conn.Write(v)
	return err
}

func (c *noopCodec) String() string {
	return "noop"
}

// NewCodec returns new noop codec
func NewCodec() Codec {
	return &noopCodec{}
}

func (c *noopCodec) Marshal(v interface{}) ([]byte, error) {
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
	case *Message:
		return ve.Body, nil
	}
	return nil, ErrInvalidMessage
}

func (c *noopCodec) Unmarshal(d []byte, v interface{}) error {
	var err error
	if v == nil {
		return nil
	}
	switch ve := v.(type) {
	case string:
		ve = string(d)
	case *string:
		*ve = string(d)
	case []byte:
		ve = d
	case *[]byte:
		*ve = d
	case *Frame:
		ve.Data = d
	case *Message:
		ve.Body = d
	default:
		err = ErrInvalidMessage
	}

	return err
}
