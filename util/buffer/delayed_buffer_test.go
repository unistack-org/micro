package buffer

import (
	"bytes"
	"testing"
	"time"
)

func TestTimedBuffer(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	b := NewDelayedBuffer(100, 300*time.Millisecond, buf)
	for range 100 {
		_, _ = b.Write([]byte(`test`))
	}
	if buf.Len() != 0 {
		t.Fatal("delayed write not worked")
	}
	time.Sleep(400 * time.Millisecond)
	if buf.Len() == 0 {
		t.Fatal("delayed write not worked")
	}
}
