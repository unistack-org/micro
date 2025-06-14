package buffer

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSeekerBuffer(t *testing.T) {
	input := []byte{'a', 'b', 'c', 'd', 'e'}
	expected := &SeekerBuffer{data: []byte{'a', 'b', 'c', 'd', 'e'}, pos: 0}
	require.Equal(t, expected, NewSeekerBuffer(input))
}

func TestSeekerBuffer_Read(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		initPos      int64
		readBuf      []byte
		expectedN    int
		expectedData []byte
		expectedErr  error
		expectedPos  int64
	}{
		{
			name:         "read with empty buffer",
			data:         []byte("hello"),
			initPos:      0,
			readBuf:      []byte{},
			expectedN:    0,
			expectedData: []byte{},
			expectedErr:  nil,
			expectedPos:  0,
		},
		{
			name:         "read with nil buffer",
			data:         []byte("hello"),
			initPos:      0,
			readBuf:      nil,
			expectedN:    0,
			expectedData: nil,
			expectedErr:  nil,
			expectedPos:  0,
		},
		{
			name:         "negative position",
			data:         []byte("hello"),
			initPos:      -1,
			readBuf:      make([]byte, 5),
			expectedN:    0,
			expectedData: make([]byte, 5),
			expectedErr:  fmt.Errorf("seeker position out of range: %d", -1),
			expectedPos:  -1,
		},
		{
			name:         "read full buffer",
			data:         []byte("hello"),
			initPos:      0,
			readBuf:      make([]byte, 5),
			expectedN:    5,
			expectedData: []byte("hello"),
			expectedErr:  nil,
			expectedPos:  5,
		},
		{
			name:         "read partial buffer",
			data:         []byte("hello"),
			initPos:      2,
			readBuf:      make([]byte, 2),
			expectedN:    2,
			expectedData: []byte("ll"),
			expectedErr:  nil,
			expectedPos:  4,
		},
		{
			name:         "read after end",
			data:         []byte("hello"),
			initPos:      5,
			readBuf:      make([]byte, 5),
			expectedN:    0,
			expectedData: make([]byte, 5),
			expectedErr:  io.EOF,
			expectedPos:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSeekerBuffer(tt.data)
			sb.pos = tt.initPos

			n, err := sb.Read(tt.readBuf)

			if tt.expectedErr != nil {
				require.Equal(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedN, n)
			require.Equal(t, tt.expectedData, tt.readBuf)
			require.Equal(t, tt.expectedPos, sb.pos)
		})
	}
}

func TestSeekerBuffer_Write(t *testing.T) {
	tests := []struct {
		name         string
		initialData  []byte
		initialPos   int64
		writeData    []byte
		expectedData []byte
		expectedN    int
	}{
		{
			name:         "write empty slice",
			initialData:  []byte("data"),
			initialPos:   0,
			writeData:    []byte{},
			expectedData: []byte("data"),
			expectedN:    0,
		},
		{
			name:         "write nil slice",
			initialData:  []byte("data"),
			initialPos:   0,
			writeData:    nil,
			expectedData: []byte("data"),
			expectedN:    0,
		},
		{
			name:         "write to empty buffer",
			initialData:  nil,
			initialPos:   0,
			writeData:    []byte("abc"),
			expectedData: []byte("abc"),
			expectedN:    3,
		},
		{
			name:         "write to existing buffer",
			initialData:  []byte("hello"),
			initialPos:   0,
			writeData:    []byte(" world"),
			expectedData: []byte("hello world"),
			expectedN:    6,
		},
		{
			name:         "write after read",
			initialData:  []byte("abc"),
			initialPos:   2,
			writeData:    []byte("XYZ"),
			expectedData: []byte("abcXYZ"),
			expectedN:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSeekerBuffer(tt.initialData)
			sb.pos = tt.initialPos

			n, err := sb.Write(tt.writeData)
			require.NoError(t, err)
			require.Equal(t, tt.expectedN, n)
			require.Equal(t, tt.expectedData, sb.data)
			require.Equal(t, tt.initialPos, sb.pos)
		})
	}
}

func TestSeekerBuffer_Seek(t *testing.T) {
	tests := []struct {
		name        string
		initialData []byte
		initialPos  int64
		offset      int64
		whence      int
		expectedPos int64
		expectedErr error
	}{
		{
			name:        "seek with invalid whence",
			initialData: []byte("abcdef"),
			initialPos:  0,
			offset:      1,
			whence:      12345,
			expectedPos: 0,
			expectedErr: fmt.Errorf("invalid whence: %d", 12345),
		},
		{
			name:        "seek negative from start",
			initialData: []byte("abcdef"),
			initialPos:  0,
			offset:      -1,
			whence:      io.SeekStart,
			expectedPos: 0,
			expectedErr: fmt.Errorf("invalid seek: resulting position %d is negative", -1),
		},
		{
			name:        "seek from start to 0",
			initialData: []byte("abcdef"),
			initialPos:  0,
			offset:      0,
			whence:      io.SeekStart,
			expectedPos: 0,
			expectedErr: nil,
		},
		{
			name:        "seek from start to 3",
			initialData: []byte("abcdef"),
			initialPos:  0,
			offset:      3,
			whence:      io.SeekStart,
			expectedPos: 3,
			expectedErr: nil,
		},
		{
			name:        "seek from end to -1 (last byte)",
			initialData: []byte("abcdef"),
			initialPos:  0,
			offset:      -1,
			whence:      io.SeekEnd,
			expectedPos: 5,
			expectedErr: nil,
		},
		{
			name:        "seek from current forward",
			initialData: []byte("abcdef"),
			initialPos:  2,
			offset:      2,
			whence:      io.SeekCurrent,
			expectedPos: 4,
			expectedErr: nil,
		},
		{
			name:        "seek from current backward",
			initialData: []byte("abcdef"),
			initialPos:  4,
			offset:      -2,
			whence:      io.SeekCurrent,
			expectedPos: 2,
			expectedErr: nil,
		},
		{
			name:        "seek to end exactly",
			initialData: []byte("abcdef"),
			initialPos:  0,
			offset:      0,
			whence:      io.SeekEnd,
			expectedPos: 6,
			expectedErr: nil,
		},
		{
			name:        "seek to out of range",
			initialData: []byte("abcdef"),
			initialPos:  0,
			offset:      2,
			whence:      io.SeekEnd,
			expectedPos: 8,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSeekerBuffer(tt.initialData)
			sb.pos = tt.initialPos

			newPos, err := sb.Seek(tt.offset, tt.whence)

			if tt.expectedErr != nil {
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPos, newPos)
				require.Equal(t, tt.expectedPos, sb.pos)
			}
		})
	}
}

func TestSeekerBuffer_Rewind(t *testing.T) {
	buf := NewSeekerBuffer([]byte("hello world"))
	buf.pos = 4

	require.NoError(t, buf.Rewind())
	require.Equal(t, []byte("hello world"), buf.data)
	require.Equal(t, int64(0), buf.pos)
}

func TestSeekerBuffer_Close(t *testing.T) {
	buf := NewSeekerBuffer([]byte("hello world"))
	buf.pos = 2

	require.NoError(t, buf.Close())
	require.Nil(t, buf.data)
	require.Equal(t, int64(0), buf.pos)
}

func TestSeekerBuffer_Reset(t *testing.T) {
	buf := NewSeekerBuffer([]byte("hello world"))
	buf.pos = 2

	buf.Reset()
	require.Nil(t, buf.data)
	require.Equal(t, int64(0), buf.pos)
}

func TestSeekerBuffer_Len(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		pos      int64
		expected int
	}{
		{
			name:     "full buffer",
			data:     []byte("abcde"),
			pos:      0,
			expected: 5,
		},
		{
			name:     "partial read",
			data:     []byte("abcde"),
			pos:      2,
			expected: 3,
		},
		{
			name:     "fully read",
			data:     []byte("abcde"),
			pos:      5,
			expected: 0,
		},
		{
			name:     "pos > len",
			data:     []byte("abcde"),
			pos:      10,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewSeekerBuffer(tt.data)
			buf.pos = tt.pos
			require.Equal(t, tt.expected, buf.Len())
		})
	}
}

func TestSeekerBuffer_Bytes(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		pos      int64
		expected []byte
	}{
		{
			name:     "start of buffer",
			data:     []byte("abcde"),
			pos:      0,
			expected: []byte("abcde"),
		},
		{
			name:     "middle of buffer",
			data:     []byte("abcde"),
			pos:      2,
			expected: []byte("cde"),
		},
		{
			name:     "end of buffer",
			data:     []byte("abcde"),
			pos:      5,
			expected: []byte{},
		},
		{
			name:     "pos beyond end",
			data:     []byte("abcde"),
			pos:      10,
			expected: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewSeekerBuffer(tt.data)
			buf.pos = tt.pos
			require.Equal(t, tt.expected, buf.Bytes())
		})
	}
}
