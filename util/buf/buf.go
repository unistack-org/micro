package buf

import (
	"bytes"
)

type buffer struct {
	*bytes.Buffer
}

// Close reset buffer contents
func (b *buffer) Close() error {
	b.Buffer.Reset()
	return nil
}

// New creates new buffer that satisfies Closer interface
func New(b *bytes.Buffer) *buffer {
	if b == nil {
		b = bytes.NewBuffer(nil)
	}
	return &buffer{b}
}
