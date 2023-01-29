package logger

import (
	"bytes"
	"context"
	"log"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	nl, ok := FromContext(NewContext(ctx, l.Fields("key", "val")))
	if !ok {
		t.Fatal("context without logger")
	}
	nl.Info(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
}

func TestFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	nl := l.Fields("key", "val")

	nl.Info(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
}

func TestFromContextWithFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	var ok bool
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	nl := l.Fields("key", "val")

	ctx = NewContext(ctx, nl)

	l, ok = FromContext(ctx)
	if !ok {
		t.Fatalf("context does not have logger")
	}

	l.Info(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
}

func TestClone(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	nl := l.Clone(WithLevel(ErrorLevel))
	if err := nl.Init(); err != nil {
		t.Fatal(err)
	}
	nl.Info(ctx, "info message")
	if len(buf.Bytes()) != 0 {
		t.Fatal("message must not be logged")
	}
	l.Info(ctx, "info message")
	if len(buf.Bytes()) == 0 {
		t.Fatal("message must be logged")
	}
}

func TestRedirectStdLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	fn := RedirectStdLogger(l, ErrorLevel)
	defer fn()
	log.Print("test")
	if !bytes.Contains(buf.Bytes(), []byte(`"level":"error","msg":"test","timestamp"`)) {
		t.Fatalf("logger error, buf %s", buf.Bytes())
	}
}

func TestStdLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	lg := NewStdLogger(l, ErrorLevel)
	lg.Print("test")
	if !bytes.Contains(buf.Bytes(), []byte(`"level":"error","msg":"test","timestamp"`)) {
		t.Fatalf("logger error, buf %s", buf.Bytes())
	}
}

func TestLogger(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(WithLevel(TraceLevel), WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Trace(ctx, "trace_msg1")
	l.Warn(ctx, "warn_msg1")
	l.Fields("error", "test").Info(ctx, "error message")
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
