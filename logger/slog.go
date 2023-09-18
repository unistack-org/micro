package logger

import (
	"context"
	"fmt"
	"go.unistack.org/micro/v4/options"
	"log"
	"sync"

	"golang.org/x/exp/slog"
)

const (
	slogName = "slog"
)

type slogLogger struct {
	slog   *slog.Logger
	fields map[string]interface{}
	opts   Options

	sync.RWMutex
}

//TODO:!!!!

func (s *slogLogger) Clone(opts ...options.Option) Logger {
	//TODO implement me
	panic("implement me")
}

func (s *slogLogger) V(level Level) bool {
	//TODO implement me
	panic("implement me")
}

func (s *slogLogger) Level(level Level) {
	//TODO implement me
	panic("implement me")
}

func (s *slogLogger) Options() Options {
	//TODO implement me
	panic("implement me")
}

func (s *slogLogger) Fields(fields ...interface{}) Logger {
	//TODO implement me
	panic("implement me")
}

func (s *slogLogger) Trace(ctx context.Context, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (s *slogLogger) Tracef(ctx context.Context, msg string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (s *slogLogger) Fatalf(ctx context.Context, msg string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (s *slogLogger) Init(opts ...options.Option) error {
	for _, o := range opts {
		if err := o(&s.opts); err != nil {
			return err
		}
	}

	if slog, ok := s.opts.Context.Value(loggerKey{}).(*slog.Logger); ok {
		s.slog = slog
		return nil
	}

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: renameTime,
		Level:       loggerToSlogLevel(s.opts.Level),
	}

	attr := fieldsToAttr(s.fields)

	handler := slog.NewJSONHandler(s.opts.Out, handleOpt).WithAttrs(attr)

	s.slog = slog.New(handler)

	return nil
}

func (s *slogLogger) Log(ctx context.Context, lvl Level, args ...any) {
	slvl := loggerToSlogLevel(lvl)

	s.RLock()
	attr := fieldsToAttr(s.fields)
	s.RUnlock()

	msg := fmt.Sprint(args...)

	if lvl == FatalLevel {
		log.Fatalln(msg, attr)
	}

	s.slog.LogAttrs(ctx, slvl, msg, attr...)
}

func (s *slogLogger) Logf(ctx context.Context, lvl Level, format string, args ...any) {
	slvl := loggerToSlogLevel(lvl)

	s.RLock()
	attr := fieldsToAttr(s.fields)
	s.RUnlock()

	msg := fmt.Sprintf(format, args...)

	if lvl == FatalLevel {
		log.Fatalln(msg, attr)
	}

	s.slog.LogAttrs(ctx, slvl, msg, attr...)
}

func (s *slogLogger) Info(ctx context.Context, args ...any) {
	s.Log(ctx, InfoLevel, args...)
}

func (s *slogLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	s.Logf(ctx, InfoLevel, format, args...)
}

func (s *slogLogger) Debug(ctx context.Context, args ...any) {
	s.Log(ctx, DebugLevel, args...)
}

func (s *slogLogger) Debugf(ctx context.Context, format string, args ...any) {
	s.Logf(ctx, DebugLevel, format, args...)
}

func (s *slogLogger) Error(ctx context.Context, args ...any) {
	s.Log(ctx, ErrorLevel, args...)
}

func (s *slogLogger) Errorf(ctx context.Context, format string, args ...any) {
	s.Logf(ctx, ErrorLevel, format, args...)
}

func (s *slogLogger) Fatal(ctx context.Context, args ...any) {
	s.Log(ctx, FatalLevel, args...)
}

func (s *slogLogger) Warn(ctx context.Context, args ...any) {
	s.Log(ctx, WarnLevel, args...)
}

func (s *slogLogger) Warnf(ctx context.Context, format string, args ...any) {
	s.Logf(ctx, WarnLevel, format, args...)
}

func (s *slogLogger) String() string {
	return slogName
}

/*
func (s *slogLogger) Fields(fields map[string]interface{}) Logger {
	nfields := make(map[string]interface{}, len(s.fields))

	s.Lock()
	for k, v := range s.fields {
		nfields[k] = v
	}
	s.Unlock()

	for k, v := range fields {
		nfields[k] = v
	}

	keys := make([]string, 0, len(nfields))
	for k := range nfields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	attr := make([]slog.Attr, 0, len(nfields))
	for _, k := range keys {
		attr = append(attr, slog.Any(k, fields[k]))
	}

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: renameTime,
		Level:       loggerToSlogLevel(s.opts.Level),
	}

	handler := slog.NewJSONHandler(s.opts.Out, handleOpt).WithAttrs(attr)

	zl := &slogLogger{
		slog:   slog.New(handler),
		opts:   s.opts,
		fields: make(map[string]interface{}),
	}

	return zl
}

*/

func NewSlogLogger(opts ...options.Option) (Logger, error) {
	l := &slogLogger{
		opts: NewOptions(opts...),
	}
	err := l.Init()
	return l, err
}

func renameTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		a.Key = "@timestamp"
	}

	return a
}

func loggerToSlogLevel(level Level) slog.Level {
	switch level {
	case TraceLevel, DebugLevel:
		return slog.LevelDebug
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel, FatalLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func fieldsToAttr(m map[string]any) []slog.Attr {
	data := make([]slog.Attr, 0, len(m))
	for k, v := range m {
		data = append(data, slog.Any(k, v))
	}

	return data
}
