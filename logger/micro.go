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

func init() {
	lvl, err := GetLevel(os.Getenv("MICRO_LOG_LEVEL"))
	if err != nil {
		lvl = InfoLevel
	}

	DefaultLogger = NewLogger(WithLevel(lvl))
}

type defaultLogger struct {
	sync.RWMutex
	opts Options
	enc  *json.Encoder
}

// Init(opts...) should only overwrite provided options
func (l *defaultLogger) Init(opts ...Option) error {
	l.Lock()
	for _, o := range opts {
		o(&l.opts)
	}
	l.enc = json.NewEncoder(l.opts.Out)
	l.Unlock()
	return nil
}

func (l *defaultLogger) String() string {
	return "micro"
}

func (l *defaultLogger) V(level Level) bool {
	l.RLock()
	ok := l.opts.Level.Enabled(level)
	l.RUnlock()
	return ok
}

func (l *defaultLogger) Fields(fields map[string]interface{}) Logger {
	l.Lock()
	l.opts.Fields = copyFields(fields)
	l.Unlock()
	return l
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
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

	fields["level"] = level.String()

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields["caller"] = fmt.Sprintf("%s:%d", logCallerfilePath(file), line)
	}

	fields["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	fields["msg"] = fmt.Sprint(args...)

	l.RLock()
	_ = l.enc.Encode(fields)
	l.RUnlock()
}

func (l *defaultLogger) Logf(ctx context.Context, level Level, msg string, args ...interface{}) {
	if !l.V(level) {
		return
	}

	l.RLock()
	fields := copyFields(l.opts.Fields)
	l.RUnlock()

	fields["level"] = level.String()

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields["caller"] = fmt.Sprintf("%s:%d", logCallerfilePath(file), line)
	}

	fields["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	if len(args) > 0 {
		fields["msg"] = fmt.Sprintf(msg, args...)
	} else {
		fields["msg"] = msg
	}
	l.RLock()
	_ = l.enc.Encode(fields)
	l.RUnlock()
}

func (l *defaultLogger) Options() Options {
	// not guard against options Context values
	l.RLock()
	opts := l.opts
	opts.Fields = copyFields(l.opts.Fields)
	l.RUnlock()
	return opts
}

// NewLogger builds a new logger based on options
func NewLogger(opts ...Option) Logger {
	l := &defaultLogger{opts: NewOptions(opts...)}
	l.enc = json.NewEncoder(l.opts.Out)
	return l
}
