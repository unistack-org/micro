package buffer

import (
	"fmt"
	"io"
)

var _ interface {
	io.ReadCloser
	io.ReadSeeker
} = (*SeekerBuffer)(nil)

// SeekerBuffer is a ReadWriteCloser that supports seeking. It's intended to
// replicate the functionality of bytes.Buffer that I use in my projects.
//
// Note that the seeking is limited to the read marker; all writes are
// append-only.
type SeekerBuffer struct {
	data []byte
	pos  int64
}

func NewSeekerBuffer(data []byte) *SeekerBuffer {
	return &SeekerBuffer{
		data: data,
	}
}

// Read reads up to len(p) bytes into p from the current read position.
func (b *SeekerBuffer) Read(p []byte) (int, error) {
	if b.pos < 0 {
		return 0, fmt.Errorf("seeker position out of range: %d", b.pos)
	}

	if b.pos >= int64(len(b.data)) {
		return 0, io.EOF
	}

	n := copy(p, b.data[b.pos:])
	b.pos += int64(n)

	return n, nil
}

// Write appends the contents of p to the end of the buffer. It does not affect the read position.
func (b *SeekerBuffer) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	b.data = append(b.data, p...)

	return len(p), nil
}

// Seek sets the offset for the next Read operation.
// The offset is interpreted according to whence:
// - io.SeekStart: relative to the beginning of the buffer
// - io.SeekCurrent: relative to the current position
// - io.SeekEnd: relative to the end of the buffer
//
// Returns an error if the resulting position is negative or if whence is invalid.
func (b *SeekerBuffer) Seek(offset int64, whence int) (int64, error) {
	var newPos int64

	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekEnd:
		newPos = int64(len(b.data)) + offset
	case io.SeekCurrent:
		newPos = b.pos + offset
	default:
		return 0, fmt.Errorf("invalid whence: %d", whence)
	}

	if newPos < 0 {
		return 0, fmt.Errorf("invalid seek: resulting position %d is negative", newPos)
	}

	b.pos = newPos
	return b.pos, nil
}

// Rewind resets the read position to 0.
func (b *SeekerBuffer) Rewind() error {
	if _, err := b.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return nil
}

// Close clears all the data out of the buffer and sets the read position to 0.
func (b *SeekerBuffer) Close() error {
	b.data = nil
	b.pos = 0
	return nil
}

// Reset clears all the data out of the buffer and sets the read position to 0.
func (b *SeekerBuffer) Reset() {
	b.data = nil
	b.pos = 0
}

// Len returns the length of data remaining to be read.
func (b *SeekerBuffer) Len() int {
	if b.pos >= int64(len(b.data)) {
		return 0
	}
	return len(b.data[b.pos:])
}

// Bytes returns the underlying bytes from the current position.
func (b *SeekerBuffer) Bytes() []byte {
	return b.data[b.pos:]
}
