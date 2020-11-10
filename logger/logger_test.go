package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	l := NewLogger(WithLevel(TraceLevel))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Trace("trace_msg1")
	l.Warn("warn_msg1")
	l.Fields(map[string]interface{}{"error": "test"}).Info("error message")
	l.Warn("first", " ", "second")
}
