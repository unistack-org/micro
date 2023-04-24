package io

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"
	"time"
)

var ErrorPattern = regexp.MustCompile(`"error"`)

func TestRedirect(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan string)

	go func() {
		buf := bytes.NewBuffer(nil)
		_, _ = io.Copy(buf, r)
		ch <- buf.String()
	}()

	if err = RedirectStderr(w); err != nil {
		t.Fatal(err)
	}

	os.Stderr.Write([]byte(`test redirect`))
	time.Sleep(1 * time.Millisecond)
	r.Close()
	str := <-ch
	if ErrorPattern.MatchString(str) {
		t.Fatal(fmt.Errorf(str))
	}
}
