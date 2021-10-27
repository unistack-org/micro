package buf // import "go.unistack.org/micro/v3/util/buf"

import (
	"bytes"
)

type Buffer struct {
	*bytes.Buffer
}

// Close reset buffer contents
func (b *Buffer) Close() error {
	b.Buffer.Reset()
	return nil
}

// New creates new buffer that satisfies Closer interface
func New(b *bytes.Buffer) *Buffer {
	if b == nil {
		b = bytes.NewBuffer(nil)
	}
	return &Buffer{b}
}
