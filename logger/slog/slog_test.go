package slog

import (
	"bytes"
	"context"
	"log"
	"testing"

	"go.unistack.org/micro/v4/logger"
)

func TestContext(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	nl, ok := logger.FromContext(logger.NewContext(ctx, l.Attrs("key", "val")))
	if !ok {
		t.Fatal("context without logger")
	}
	nl.Info(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
}

func TestAttrs(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	nl := l.Attrs("key", "val")

	nl.Info(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
}

func TestFromContextWithFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	var ok bool
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	nl := l.Attrs("key", "val")

	ctx = logger.NewContext(ctx, nl)

	l, ok = logger.FromContext(ctx)
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
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	nl := l.Clone(logger.WithLevel(logger.ErrorLevel))
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
	l := NewLogger(logger.WithLevel(logger.ErrorLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	fn := logger.RedirectStdLogger(l, logger.ErrorLevel)
	defer fn()
	log.Print("test")
	if !(bytes.Contains(buf.Bytes(), []byte(`"level":"error"`)) && bytes.Contains(buf.Bytes(), []byte(`"msg":"test"`))) {
		t.Fatalf("logger error, buf %s", buf.Bytes())
	}
}

func TestStdLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	lg := logger.NewStdLogger(l, logger.ErrorLevel)
	lg.Print("test")
	if !(bytes.Contains(buf.Bytes(), []byte(`"level":"error"`)) && bytes.Contains(buf.Bytes(), []byte(`"msg":"test"`))) {
		t.Fatalf("logger error, buf %s", buf.Bytes())
	}
}

func TestLogger(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Trace(ctx, "trace_msg1")
	l.Warn(ctx, "warn_msg1")
	l.Attrs("error", "test").Info(ctx, "error message")
	l.Warn(ctx, "first second")

	if !(bytes.Contains(buf.Bytes(), []byte(`"level":"trace"`)) && bytes.Contains(buf.Bytes(), []byte(`"msg":"trace_msg1"`))) {
		t.Fatalf("logger tracer, buf %s", buf.Bytes())
	}
	if !(bytes.Contains(buf.Bytes(), []byte(`"level":"warn"`)) && bytes.Contains(buf.Bytes(), []byte(`"msg":"warn_msg1"`))) {
		t.Fatalf("logger warn, buf %s", buf.Bytes())
	}
	if !(bytes.Contains(buf.Bytes(), []byte(`"level":"info"`)) && bytes.Contains(buf.Bytes(), []byte(`"msg":"error message","error":"test"`))) {
		t.Fatalf("logger info, buf %s", buf.Bytes())
	}
	if !(bytes.Contains(buf.Bytes(), []byte(`"level":"warn"`)) && bytes.Contains(buf.Bytes(), []byte(`"msg":"first second"`))) {
		t.Fatalf("logger warn, buf %s", buf.Bytes())
	}
}
