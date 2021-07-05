package logger

import (
	"bytes"
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Trace(ctx, "trace_msg1")
	l.Warn(ctx, "warn_msg1")
	l.Fields(map[string]interface{}{"error": "test"}).Info(ctx, "error message")
	l.Warn(ctx, "first", " ", "second")
	if !bytes.Contains(buf.Bytes(), []byte(`"level":"trace","msg":"trace_msg1"`)) {
		t.Fatalf("logger error, buf %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"warn","msg":"warn_msg1"`)) {
		t.Fatalf("logger error, buf %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"error":"test","level":"info","msg":"error message"`)) {
		t.Fatalf("logger error, buf %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"level":"warn","msg":"first second"`)) {
		t.Fatalf("logger error, buf %s", buf.Bytes())
	}
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

func TestOmitLoggerWrapper(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewOmitLogger(NewLogger(WithLevel(TraceLevel), WithOutput(buf)))
	if err := l.Init(); err != nil {
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
