package slog

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"time"

	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/tracer"
)

var (
	traceValue = slog.StringValue("trace")
	debugValue = slog.StringValue("debug")
	infoValue  = slog.StringValue("info")
	warnValue  = slog.StringValue("warn")
	errorValue = slog.StringValue("error")
	fatalValue = slog.StringValue("fatal")
)

func (s *slogLogger) renameAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.SourceKey:
		source := a.Value.Any().(*slog.Source)
		a.Value = slog.StringValue(source.File + ":" + strconv.Itoa(source.Line))
		a.Key = s.opts.SourceKey
	case slog.TimeKey:
		a.Key = s.opts.TimeKey
	case slog.MessageKey:
		a.Key = s.opts.MessageKey
	case slog.LevelKey:
		level := a.Value.Any().(slog.Level)
		lvl := slogToLoggerLevel(level)
		a.Key = s.opts.LevelKey
		switch {
		case lvl < logger.DebugLevel:
			a.Value = traceValue
		case lvl < logger.InfoLevel:
			a.Value = debugValue
		case lvl < logger.WarnLevel:
			a.Value = infoValue
		case lvl < logger.ErrorLevel:
			a.Value = warnValue
		case lvl < logger.FatalLevel:
			a.Value = errorValue
		case lvl >= logger.FatalLevel:
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
	opts    logger.Options
}

func (s *slogLogger) Clone(opts ...options.Option) logger.Logger {
	options := s.opts

	for _, o := range opts {
		o(&options)
	}

	l := &slogLogger{
		opts: options,
	}

	/*
		if slog, ok := s.opts.Context.Value(loggerKey{}).(*slog.Logger); ok {
			l.slog = slog
			return nil
		}
	*/

	l.leveler = new(slog.LevelVar)
	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: s.renameAttr,
		Level:       l.leveler,
		AddSource:   true,
	}
	l.leveler.Set(loggerToSlogLevel(l.opts.Level))
	handler := slog.NewJSONHandler(options.Out, handleOpt)
	l.slog = slog.New(handler).With(options.Attrs...)

	return l
}

func (s *slogLogger) V(level logger.Level) bool {
	return s.opts.Level.Enabled(level)
}

func (s *slogLogger) Level(level logger.Level) {
	s.leveler.Set(loggerToSlogLevel(level))
}

func (s *slogLogger) Options() logger.Options {
	return s.opts
}

func (s *slogLogger) Attrs(attrs ...interface{}) logger.Logger {
	nl := &slogLogger{opts: s.opts}
	nl.leveler = new(slog.LevelVar)
	nl.leveler.Set(s.leveler.Level())

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: nl.renameAttr,
		Level:       nl.leveler,
		AddSource:   true,
	}

	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	nl.slog = slog.New(handler).With(attrs...)

	return nl
}

func (s *slogLogger) Init(opts ...options.Option) error {
	if len(s.opts.ContextAttrFuncs) == 0 {
		s.opts.ContextAttrFuncs = logger.DefaultContextAttrFuncs
	}
	for _, o := range opts {
		if err := o(&s.opts); err != nil {
			return err
		}
	}

	/*
		if slog, ok := s.opts.Context.Value(loggerKey{}).(*slog.Logger); ok {
			s.slog = slog
			return nil
		}
	*/

	s.leveler = new(slog.LevelVar)
	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: s.renameAttr,
		Level:       s.leveler,
		AddSource:   true,
	}
	s.leveler.Set(loggerToSlogLevel(s.opts.Level))
	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	s.slog = slog.New(handler).With(s.opts.Attrs...)

	return nil
}

func (s *slogLogger) Log(ctx context.Context, lvl logger.Level, msg string, attrs ...interface{}) {
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
	if !s.V(logger.InfoLevel) {
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
	if !s.V(logger.DebugLevel) {
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
	if !s.V(logger.TraceLevel) {
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
	if !s.V(logger.ErrorLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}
	r.Add(attrs...)
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "error" {
			if span, ok := tracer.SpanFromContext(ctx); ok {
				span.SetStatus(tracer.SpanStatusError, a.Value.String())
				return false
			}
		}
		return true
	})
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Fatal(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.FatalLevel) {
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
	if !s.V(logger.WarnLevel) {
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

func (s *slogLogger) Name() string {
	return s.opts.Name
}

func (s *slogLogger) String() string {
	return "slog"
}

func NewLogger(opts ...options.Option) logger.Logger {
	l := &slogLogger{
		opts: logger.NewOptions(opts...),
	}
	return l
}

func loggerToSlogLevel(level logger.Level) slog.Level {
	switch level {
	case logger.DebugLevel:
		return slog.LevelDebug
	case logger.WarnLevel:
		return slog.LevelWarn
	case logger.ErrorLevel:
		return slog.LevelError
	case logger.TraceLevel:
		return slog.LevelDebug - 1
	case logger.FatalLevel:
		return slog.LevelError + 1
	default:
		return slog.LevelInfo
	}
}

func slogToLoggerLevel(level slog.Level) logger.Level {
	switch level {
	case slog.LevelDebug:
		return logger.DebugLevel
	case slog.LevelWarn:
		return logger.WarnLevel
	case slog.LevelError:
		return logger.ErrorLevel
	case slog.LevelDebug - 1:
		return logger.TraceLevel
	case slog.LevelError + 1:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}
