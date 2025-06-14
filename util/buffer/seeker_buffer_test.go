package buffer

import (
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
		name     string
		data     []byte
		initPos  int64
		readBuf  []byte
		wantN    int
		wantData []byte
		wantErr  error
	}{
		{
			name:     "read with empty buffer",
			data:     []byte("hello"),
			initPos:  0,
			readBuf:  []byte{},
			wantN:    0,
			wantData: []byte{},
			wantErr:  nil,
		},
		{
			name:     "read with nil buffer",
			data:     []byte("hello"),
			initPos:  0,
			readBuf:  nil,
			wantN:    0,
			wantData: nil,
			wantErr:  nil,
		},
		{
			name:     "negative position",
			data:     []byte("hello"),
			initPos:  -1,
			readBuf:  make([]byte, 5),
			wantN:    0,
			wantData: make([]byte, 5),
			wantErr:  ErrNegativePosition,
		},
		{
			name:     "read full buffer",
			data:     []byte("hello"),
			initPos:  0,
			readBuf:  make([]byte, 5),
			wantN:    5,
			wantData: []byte("hello"),
			wantErr:  nil,
		},
		{
			name:     "read partial buffer",
			data:     []byte("hello"),
			initPos:  2,
			readBuf:  make([]byte, 2),
			wantN:    2,
			wantData: []byte("ll"),
			wantErr:  nil,
		},
		{
			name:     "read after end",
			data:     []byte("hello"),
			initPos:  5,
			readBuf:  make([]byte, 5),
			wantN:    0,
			wantData: make([]byte, 5),
			wantErr:  io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSeekerBuffer(tt.data)
			sb.pos = tt.initPos

			n, err := sb.Read(tt.readBuf)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantN, n)
			require.Equal(t, tt.wantData, tt.readBuf)
		})
	}
}
