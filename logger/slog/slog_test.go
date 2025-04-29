package slog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/util/buffer"
)

// always first to have proper check
func TestStacktrace(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.DebugLevel), logger.WithOutput(buf),
		WithHandlerFunc(slog.NewTextHandler),
		logger.WithAddStacktrace(true),
	)
	if err := l.Init(logger.WithFields("key1", "val1")); err != nil {
		t.Fatal(err)
	}

	l.Error(ctx, "msg1", errors.New("err"))

	if !bytes.Contains(buf.Bytes(), []byte(`slog_test.go:32`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestDelayedBuffer(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	dbuf := buffer.NewDelayedBuffer(100, 100*time.Millisecond, buf)
	l := NewLogger(logger.WithLevel(logger.ErrorLevel), logger.WithOutput(dbuf),
		WithHandlerFunc(slog.NewTextHandler),
		logger.WithAddStacktrace(true),
	)
	if err := l.Init(logger.WithFields("key1", "val1")); err != nil {
		t.Fatal(err)
	}

	l.Error(ctx, "msg1", errors.New("err"))
	time.Sleep(120 * time.Millisecond)
	if !bytes.Contains(buf.Bytes(), []byte(`key1=val1`)) {
		t.Fatalf("logger delayed buffer not works, buf contains: %s", buf.Bytes())
	}
}

func TestTime(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.ErrorLevel), logger.WithOutput(buf),
		WithHandlerFunc(slog.NewTextHandler),
		logger.WithAddStacktrace(true),
		logger.WithTimeFunc(func() time.Time {
			return time.Unix(0, 0)
		}),
	)
	if err := l.Init(logger.WithFields("key1", "val1")); err != nil {
		t.Fatal(err)
	}

	l.Error(ctx, "msg1", errors.New("err"))

	if !bytes.Contains(buf.Bytes(), []byte(`timestamp=1970-01-01T03:00:00.000000000+03:00`)) &&
		!bytes.Contains(buf.Bytes(), []byte(`timestamp=1970-01-01T00:00:00.000000000Z`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestWithFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.InfoLevel), logger.WithOutput(buf),
		WithHandlerFunc(slog.NewTextHandler),
		logger.WithDedupKeys(true),
	)
	if err := l.Init(logger.WithFields("key1", "val1")); err != nil {
		t.Fatal(err)
	}

	l.Info(ctx, "msg1")

	l = l.Fields("key1", "val2")

	l.Info(ctx, "msg2")

	if !bytes.Contains(buf.Bytes(), []byte(`msg=msg2 key1=val1`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestWithDedupKeysWithAddFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.InfoLevel), logger.WithOutput(buf),
		WithHandlerFunc(slog.NewTextHandler),
		logger.WithDedupKeys(true),
	)
	if err := l.Init(logger.WithFields("key1", "val1")); err != nil {
		t.Fatal(err)
	}

	l.Info(ctx, "msg1")

	if err := l.Init(logger.WithAddFields("key2", "val2")); err != nil {
		t.Fatal(err)
	}

	l.Info(ctx, "msg2")

	if err := l.Init(logger.WithAddFields("key2", "val3", "key1", "val4")); err != nil {
		t.Fatal(err)
	}

	l.Info(ctx, "msg3")

	if !bytes.Contains(buf.Bytes(), []byte(`msg=msg3 key1=val4 key2=val3`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestWithHandlerFunc(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.InfoLevel), logger.WithOutput(buf),
		WithHandlerFunc(slog.NewTextHandler),
	)
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l.Info(ctx, "msg1")

	if !bytes.Contains(buf.Bytes(), []byte(`msg=msg1`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestWithAddFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.InfoLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l.Info(ctx, "msg1")

	if err := l.Init(logger.WithAddFields("key1", "val1")); err != nil {
		t.Fatal(err)
	}
	l.Info(ctx, "msg2")

	if err := l.Init(logger.WithAddFields("key2", "val2")); err != nil {
		t.Fatal(err)
	}
	l.Info(ctx, "msg3")

	if !bytes.Contains(buf.Bytes(), []byte(`"key1"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"key2"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestMultipleFieldsWithLevel(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.InfoLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l = l.Fields("key", "val")

	l.Info(ctx, "msg1")
	nl := l.Clone(logger.WithLevel(logger.DebugLevel))
	nl.Debug(ctx, "msg2")
	l.Debug(ctx, "msg3")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"msg1"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"msg2"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
	if bytes.Contains(buf.Bytes(), []byte(`"msg3"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestMultipleFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.InfoLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l = l.Fields("key", "val")

	l = l.Fields("key1", "val1")

	l.Info(ctx, "msg")

	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"key1":"val1"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestError(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.ErrorLevel), logger.WithOutput(buf), logger.WithAddStacktrace(true))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l.Error(ctx, "message", fmt.Errorf("error message"))
	if !bytes.Contains(buf.Bytes(), []byte(`"stacktrace":"`)) {
		t.Fatalf("logger stacktrace not works, buf contains: %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"error":"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestErrorf(t *testing.T) {
	ctx := context.TODO()

	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.ErrorLevel), logger.WithOutput(buf), logger.WithAddStacktrace(true))
	if err := l.Init(logger.WithContextAttrFuncs(func(_ context.Context) []interface{} {
		return nil
	})); err != nil {
		t.Fatal(err)
	}

	l.Log(ctx, logger.ErrorLevel, "message", errors.New("error msg"))

	l.Log(ctx, logger.ErrorLevel, "", errors.New("error msg"))
	if !bytes.Contains(buf.Bytes(), []byte(`"error":"error msg"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}

	if !bytes.Contains(buf.Bytes(), []byte(`"stacktrace":"`)) {
		t.Fatalf("logger stacktrace not works, buf contains: %s", buf.Bytes())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"error":"`)) {
		t.Fatalf("logger error not works, buf contains: %s", buf.Bytes())
	}
}

func TestContext(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	nl, ok := logger.FromContext(logger.NewContext(ctx, l.Fields("key", "val")))
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
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
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
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	nl := l.Fields("key", "val")

	ctx = logger.NewContext(ctx, nl)

	l, ok = logger.FromContext(ctx)
	if !ok {
		t.Fatalf("context does not have logger")
	}

	l.Info(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}

	l.Info(ctx, "test", "uncorrected number attributes")
	if !bytes.Contains(buf.Bytes(), []byte(`"!BADKEY":"uncorrected number attributes"`)) {
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
	l.Fields("error", "test").Info(ctx, "error message")
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

func Test_WithContextAttrFunc(t *testing.T) {
	loggerContextAttrFuncs := []logger.ContextAttrFunc{
		func(ctx context.Context) []interface{} {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				return nil
			}
			attrs := make([]interface{}, 0, 10)
			for k, v := range md {
				key := strings.ToLower(k)
				switch key {
				case "x-request-id", "phone", "external-Id", "source-service", "x-app-install-id", "client-id", "client-ip":
					attrs = append(attrs, key, v[0])
				}
			}
			return attrs
		},
	}

	logger.DefaultContextAttrFuncs = append(logger.DefaultContextAttrFuncs, loggerContextAttrFuncs...)

	ctx := context.TODO()
	ctx = metadata.AppendOutgoingContext(ctx, "X-Request-Id", uuid.New().String(),
		"Source-Service", "Test-System")

	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l.Info(ctx, "test message")
	if !(bytes.Contains(buf.Bytes(), []byte(`"level":"info"`)) && bytes.Contains(buf.Bytes(), []byte(`"msg":"test message"`))) {
		t.Fatalf("logger info, buf %s", buf.Bytes())
	}
	if !(bytes.Contains(buf.Bytes(), []byte(`"x-request-id":"`))) {
		t.Fatalf("logger info, buf %s", buf.Bytes())
	}
	if !(bytes.Contains(buf.Bytes(), []byte(`"source-service":"Test-System"`))) {
		t.Fatalf("logger info, buf %s", buf.Bytes())
	}
	buf.Reset()
	omd, _ := metadata.FromOutgoingContext(ctx)
	l.Info(ctx, "test message1")
	omd.Set("Source-Service", "Test-System2")
	l.Info(ctx, "test message2")

	// t.Logf("xxx %s", buf.Bytes())
}
