package logger

import (
	"context"
	"fmt"
	"os"

	"go.unistack.org/micro/v4/options"

	"golang.org/x/exp/slog"
)

var (
	traceValue = slog.StringValue("trace")
	debugValue = slog.StringValue("debug")
	infoValue  = slog.StringValue("info")
	warnValue  = slog.StringValue("warn")
	errorValue = slog.StringValue("error")
	fatalValue = slog.StringValue("fatal")
)

var renameAttr = func(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		a.Key = "timestamp"
	case slog.LevelKey:
		level := a.Value.Any().(slog.Level)
		lvl := slogToLoggerLevel(level)
		switch {
		case lvl < DebugLevel:
			a.Value = traceValue
		case lvl < InfoLevel:
			a.Value = debugValue
		case lvl < WarnLevel:
			a.Value = infoValue
		case lvl < ErrorLevel:
			a.Value = warnValue
		case lvl < FatalLevel:
			a.Value = errorValue
		case lvl >= FatalLevel:
			a.Value = fatalValue
		default:
			a.Value = infoValue
		}
	}

	return a
}

type slogLogger struct {
	slog    *slog.Logger
	leveler *slog.LevelVar
	opts    Options
}

func (s *slogLogger) Clone(opts ...options.Option) Logger {
	options := s.opts

	for _, o := range opts {
		o(&options)
	}

	l := &slogLogger{
		opts: options,
	}

	if slog, ok := s.opts.Context.Value(loggerKey{}).(*slog.Logger); ok {
		l.slog = slog
		return nil
	}

	l.leveler = new(slog.LevelVar)
	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: renameAttr,
		Level:       l.leveler,
	}
	l.leveler.Set(loggerToSlogLevel(l.opts.Level))
	handler := slog.NewJSONHandler(options.Out, handleOpt)
	l.slog = slog.New(handler).With(options.Fields...)

	return l
}

func (s *slogLogger) V(level Level) bool {
	return s.opts.Level.Enabled(level)
}

func (s *slogLogger) Level(level Level) {
	s.leveler.Set(loggerToSlogLevel(level))
}

func (s *slogLogger) Options() Options {
	return s.opts
}

func (s *slogLogger) Fields(fields ...interface{}) Logger {
	nl := &slogLogger{opts: s.opts}
	nl.leveler = new(slog.LevelVar)
	nl.leveler.Set(s.leveler.Level())

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: renameAttr,
		Level:       s.leveler,
	}

	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	nl.slog = slog.New(handler).With(fields...)

	return nl
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

	s.leveler = new(slog.LevelVar)
	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: renameAttr,
		Level:       s.leveler,
	}
	s.leveler.Set(loggerToSlogLevel(s.opts.Level))
	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	s.slog = slog.New(handler).With(s.opts.Fields...)

	return nil
}

func (s *slogLogger) Log(ctx context.Context, lvl Level, args ...any) {
	if !s.V(lvl) {
		return
	}
	slvl := loggerToSlogLevel(lvl)
	msg := fmt.Sprint(args...)
	s.slog.Log(ctx, slvl, msg)
}

func (s *slogLogger) Logf(ctx context.Context, lvl Level, format string, args ...any) {
	if !s.V(lvl) {
		return
	}
	slvl := loggerToSlogLevel(lvl)
	s.slog.Log(ctx, slvl, format, args...)
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

func (s *slogLogger) Trace(ctx context.Context, args ...interface{}) {
	s.Log(ctx, TraceLevel, args...)
}

func (s *slogLogger) Tracef(ctx context.Context, msg string, args ...interface{}) {
	s.Logf(ctx, TraceLevel, msg, args...)
}

func (s *slogLogger) Errorf(ctx context.Context, format string, args ...any) {
	s.Logf(ctx, ErrorLevel, format, args...)
}

func (s *slogLogger) Fatal(ctx context.Context, args ...any) {
	s.Log(ctx, FatalLevel, args...)
}

func (s *slogLogger) Fatalf(ctx context.Context, msg string, args ...interface{}) {
	s.Logf(ctx, FatalLevel, msg, args...)
	os.Exit(1)
}

func (s *slogLogger) Warn(ctx context.Context, args ...any) {
	s.Log(ctx, WarnLevel, args...)
}

func (s *slogLogger) Warnf(ctx context.Context, format string, args ...any) {
	s.Logf(ctx, WarnLevel, format, args...)
}

func (s *slogLogger) String() string {
	return "slog"
}

func NewLogger(opts ...options.Option) Logger {
	l := &slogLogger{
		opts: NewOptions(opts...),
	}
	return l
}

func loggerToSlogLevel(level Level) slog.Level {
	switch level {
	case DebugLevel:
		return slog.LevelDebug
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	case TraceLevel:
		return slog.LevelDebug - 1
	case FatalLevel:
		return slog.LevelError + 1
	default:
		return slog.LevelInfo
	}
}

func slogToLoggerLevel(level slog.Level) Level {
	switch level {
	case slog.LevelDebug:
		return DebugLevel
	case slog.LevelWarn:
		return WarnLevel
	case slog.LevelError:
		return ErrorLevel
	case slog.LevelDebug - 1:
		return TraceLevel
	case slog.LevelError + 1:
		return FatalLevel
	default:
		return InfoLevel
	}
}
