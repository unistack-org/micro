package logger

import (
	"bytes"
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.TODO()
	l := NewLogger(WithLevel(TraceLevel))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Trace(ctx, "trace_msg1")
	l.Warn(ctx, "warn_msg1")
	l.Fields(map[string]interface{}{"error": "test"}).Info(ctx, "error message")
	l.Warn(ctx, "first", " ", "second")
}

func TestLoggerWrapper(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(WrapLogger(NewOmitWrapper())); err != nil {
		t.Fatal(err)
	}
	type secret struct {
		Name  string
		Passw string `logger:"omit"`
	}
	s := &secret{Name: "name", Passw: "secret"}
	l.Errorf(ctx, "test %#+v", s)
	if !bytes.Contains(buf.Bytes(), []byte(`logger.secret{Name:\"name\", Passw:\"\"}"`)) {
		t.Fatalf("omit not works, struct: %v, output: %s", s, buf.Bytes())
	}
}
