package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"time"

	"go.unistack.org/micro/v4/options"
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
	case slog.SourceKey:
		source := a.Value.Any().(*slog.Source)
		a.Value = slog.StringValue(source.File + ":" + strconv.Itoa(source.Line))
		a.Key = "caller"
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
		AddSource:   true,
	}
	l.leveler.Set(loggerToSlogLevel(l.opts.Level))
	handler := slog.NewJSONHandler(options.Out, handleOpt)
	l.slog = slog.New(handler).With(options.Attrs...)

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

func (s *slogLogger) Attrs(attrs ...interface{}) Logger {
	nl := &slogLogger{opts: s.opts}
	nl.leveler = new(slog.LevelVar)
	nl.leveler.Set(s.leveler.Level())

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: renameAttr,
		Level:       s.leveler,
		AddSource:   true,
	}

	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	nl.slog = slog.New(handler).With(attrs...)

	return nl
}

func (s *slogLogger) Init(opts ...options.Option) error {
	if len(s.opts.ContextAttrFuncs) == 0 {
		s.opts.ContextAttrFuncs = DefaultContextAttrFuncs
	}
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
		AddSource:   true,
	}
	s.leveler.Set(loggerToSlogLevel(s.opts.Level))
	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	s.slog = slog.New(handler).With(s.opts.Attrs...)

	return nil
}

func (s *slogLogger) Log(ctx context.Context, lvl Level, msg string, attrs ...interface{}) {
	if !s.V(lvl) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), loggerToSlogLevel(lvl), msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Info(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Debug(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Trace(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelDebug-1, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Error(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Fatal(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError+1, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
	os.Exit(1)
}

func (s *slogLogger) Warn(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelWarn, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
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
