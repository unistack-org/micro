package buffer

import (
	"fmt"
	"strings"
	"testing"
)

func noErrorT(t *testing.T, err error) {
	if nil != err {
		t.Fatalf("%s", err)
	}
}

func boolT(t *testing.T, cond bool, s ...string) {
	if !cond {
		what := strings.Join(s, ", ")
		if len(what) > 0 {
			what = ": " + what
		}
		t.Fatalf("assert.Bool failed%s", what)
	}
}

func TestSeeking(t *testing.T) {
	partA := []byte("hello, ")
	partB := []byte("world!")

	buf := NewSeekerBuffer(partA)

	boolT(t, buf.Len() == len(partA), fmt.Sprintf("on init: have length %d, want length %d", buf.Len(), len(partA)))

	b := make([]byte, 32)

	n, err := buf.Read(b)
	noErrorT(t, err)
	boolT(t, buf.Len() == 0, fmt.Sprintf("after reading 1: have length %d, want length 0", buf.Len()))
	boolT(t, n == len(partA), fmt.Sprintf("after reading 2: have length %d, want length %d", n, len(partA)))

	n, err = buf.Write(partB)
	noErrorT(t, err)
	boolT(t, n == len(partB), fmt.Sprintf("after writing: have length %d, want length %d", n, len(partB)))

	n, err = buf.Read(b)
	noErrorT(t, err)
	boolT(t, buf.Len() == 0, fmt.Sprintf("after rereading 1: have length %d, want length 0", buf.Len()))
	boolT(t, n == len(partB), fmt.Sprintf("after rereading 2: have length %d, want length %d", n, len(partB)))

	partsLen := len(partA) + len(partB)
	_ = buf.Rewind()
	boolT(t, buf.Len() == partsLen, fmt.Sprintf("after rewinding: have length %d, want length %d", buf.Len(), partsLen))

	buf.Close()
	boolT(t, buf.Len() == 0, fmt.Sprintf("after closing, have length %d, want length 0", buf.Len()))
}
