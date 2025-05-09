package buffer

import "io"

var _ interface {
	io.ReadCloser
	io.ReadSeeker
} = (*SeekerBuffer)(nil)

// Buffer is a ReadWriteCloser that supports seeking. It's intended to
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

func (b *SeekerBuffer) Read(p []byte) (int, error) {
	if b.pos >= int64(len(b.data)) {
		return 0, io.EOF
	}

	n := copy(p, b.data[b.pos:])
	b.pos += int64(n)
	return n, nil
}

func (b *SeekerBuffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

// Seek sets the read pointer to pos.
func (b *SeekerBuffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.pos = offset
	case io.SeekEnd:
		b.pos = int64(len(b.data)) + offset
	case io.SeekCurrent:
		b.pos += offset
	}

	return b.pos, nil
}

// Rewind resets the read pointer to 0.
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

// Len returns the length of data remaining to be read.
func (b *SeekerBuffer) Len() int {
	return len(b.data[b.pos:])
}

// Bytes returns the underlying bytes from the current position.
func (b *SeekerBuffer) Bytes() []byte {
	return b.data[b.pos:]
}
