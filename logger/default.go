package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type defaultLogger struct {
	enc  *json.Encoder
	opts Options
	sync.RWMutex
}

// Init(opts...) should only overwrite provided options
func (l *defaultLogger) Init(opts ...Option) error {
	l.Lock()
	for _, o := range opts {
		o(&l.opts)
	}
	l.enc = json.NewEncoder(l.opts.Out)
	// wrap the Log func
	l.Unlock()
	return nil
}

func (l *defaultLogger) String() string {
	return "micro"
}

func (l *defaultLogger) Clone(opts ...Option) Logger {
	newopts := NewOptions(opts...)
	oldopts := l.opts
	for _, o := range opts {
		o(&newopts)
		o(&oldopts)
	}

	l.Lock()
	cl := &defaultLogger{opts: oldopts, enc: json.NewEncoder(l.opts.Out)}
	l.Unlock()

	return cl
}

func (l *defaultLogger) V(level Level) bool {
	l.RLock()
	ok := l.opts.Level.Enabled(level)
	l.RUnlock()
	return ok
}

func (l *defaultLogger) Level(level Level) {
	l.Lock()
	l.opts.Level = level
	l.Unlock()
}

func (l *defaultLogger) Fields(fields ...interface{}) Logger {
	l.RLock()
	nl := &defaultLogger{opts: l.opts, enc: l.enc}
	if len(fields) == 0 {
		l.RUnlock()
		return nl
	} else if len(fields)%2 != 0 {
		fields = fields[:len(fields)-1]
	}
	nl.opts.Fields = copyFields(l.opts.Fields)
	nl.opts.Fields = append(nl.opts.Fields, fields...)
	l.RUnlock()
	return nl
}

func copyFields(src []interface{}) []interface{} {
	dst := make([]interface{}, len(src))
	copy(dst, src)
	return dst
}

// logCallerfilePath returns a package/file:line description of the caller,
// preserving only the leaf directory name and file name.
func logCallerfilePath(loggingFilePath string) string {
	// To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	idx := strings.LastIndexByte(loggingFilePath, '/')
	if idx == -1 {
		return loggingFilePath
	}
	idx = strings.LastIndexByte(loggingFilePath[:idx], '/')
	if idx == -1 {
		return loggingFilePath
	}
	return loggingFilePath[idx+1:]
}

func (l *defaultLogger) Info(ctx context.Context, args ...interface{}) {
	l.Log(ctx, InfoLevel, args...)
}

func (l *defaultLogger) Error(ctx context.Context, args ...interface{}) {
	l.Log(ctx, ErrorLevel, args...)
}

func (l *defaultLogger) Debug(ctx context.Context, args ...interface{}) {
	l.Log(ctx, DebugLevel, args...)
}

func (l *defaultLogger) Warn(ctx context.Context, args ...interface{}) {
	l.Log(ctx, WarnLevel, args...)
}

func (l *defaultLogger) Trace(ctx context.Context, args ...interface{}) {
	l.Log(ctx, TraceLevel, args...)
}

func (l *defaultLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.Log(ctx, FatalLevel, args...)
	os.Exit(1)
}

func (l *defaultLogger) Infof(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, InfoLevel, msg, args...)
}

func (l *defaultLogger) Errorf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, ErrorLevel, msg, args...)
}

func (l *defaultLogger) Debugf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, DebugLevel, msg, args...)
}

func (l *defaultLogger) Warnf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, WarnLevel, msg, args...)
}

func (l *defaultLogger) Tracef(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, TraceLevel, msg, args...)
}

func (l *defaultLogger) Fatalf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, FatalLevel, msg, args...)
	os.Exit(1)
}

func (l *defaultLogger) Log(ctx context.Context, level Level, args ...interface{}) {
	if !l.V(level) {
		return
	}

	l.RLock()
	fields := copyFields(l.opts.Fields)
	l.RUnlock()

	fields = append(fields, "level", level.String())

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields = append(fields, "caller", fmt.Sprintf("%s:%d", logCallerfilePath(file), line))
	}
	fields = append(fields, "timestamp", time.Now().Format("2006-01-02 15:04:05"))

	if len(args) > 0 {
		fields = append(fields, "msg", fmt.Sprint(args...))
	}

	out := make(map[string]interface{}, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		out[fields[i].(string)] = fields[i+1]
	}
	l.RLock()
	_ = l.enc.Encode(out)
	l.RUnlock()
}

func (l *defaultLogger) Logf(ctx context.Context, level Level, msg string, args ...interface{}) {
	if !l.V(level) {
		return
	}

	l.RLock()
	fields := copyFields(l.opts.Fields)
	l.RUnlock()

	fields = append(fields, "level", level.String())

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields = append(fields, "caller", fmt.Sprintf("%s:%d", logCallerfilePath(file), line))
	}

	fields = append(fields, "timestamp", time.Now().Format("2006-01-02 15:04:05"))
	if len(args) > 0 {
		fields = append(fields, "msg", fmt.Sprintf(msg, args...))
	} else if msg != "" {
		fields = append(fields, "msg", msg)
	}

	out := make(map[string]interface{}, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		out[fields[i].(string)] = fields[i+1]
	}
	l.RLock()
	_ = l.enc.Encode(out)
	l.RUnlock()
}

func (l *defaultLogger) Options() Options {
	return l.opts
}

// NewLogger builds a new logger based on options
func NewLogger(opts ...Option) Logger {
	l := &defaultLogger{
		opts: NewOptions(opts...),
	}
	l.enc = json.NewEncoder(l.opts.Out)
	return l
}
