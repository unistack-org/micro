package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/tracer"
)

var (
	DefaultSourceKey  string = slog.SourceKey
	DefaultTimeKey    string = slog.TimeKey
	DefaultMessageKey string = slog.MessageKey
	DefaultLevelKey   string = slog.LevelKey
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
		a.Key = s.sourceKey
	case slog.TimeKey:
		a.Key = s.timeKey
	case slog.MessageKey:
		a.Key = s.messageKey
	case slog.LevelKey:
		level := a.Value.Any().(slog.Level)
		lvl := slogToLoggerLevel(level)
		a.Key = s.levelKey
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
	slog       *slog.Logger
	leveler    *slog.LevelVar
	levelKey   string
	messageKey string
	sourceKey  string
	timeKey    string
	opts       logger.Options
	mu         sync.RWMutex
}

func (s *slogLogger) Clone(opts ...logger.Option) logger.Logger {
	s.mu.RLock()
	options := s.opts

	for _, o := range opts {
		o(&options)
	}

	l := &slogLogger{
		opts:       options,
		levelKey:   s.levelKey,
		messageKey: s.messageKey,
		sourceKey:  s.sourceKey,
		timeKey:    s.timeKey,
	}

	if v, ok := l.opts.Context.Value(levelKey{}).(string); ok && v != "" {
		l.levelKey = v
	}
	if v, ok := l.opts.Context.Value(messageKey{}).(string); ok && v != "" {
		l.messageKey = v
	}
	if v, ok := l.opts.Context.Value(sourceKey{}).(string); ok && v != "" {
		l.sourceKey = v
	}
	if v, ok := l.opts.Context.Value(timeKey{}).(string); ok && v != "" {
		l.timeKey = v
	}

	l.leveler = new(slog.LevelVar)
	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: s.renameAttr,
		Level:       l.leveler,
		AddSource:   true,
	}
	l.leveler.Set(loggerToSlogLevel(l.opts.Level))
	handler := slog.NewJSONHandler(options.Out, handleOpt)
	l.slog = slog.New(handler).With(options.Fields...)

	s.mu.RUnlock()

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

func (s *slogLogger) Fields(attrs ...interface{}) logger.Logger {
	s.mu.RLock()
	nl := &slogLogger{
		opts:       s.opts,
		levelKey:   s.levelKey,
		messageKey: s.messageKey,
		sourceKey:  s.sourceKey,
		timeKey:    s.timeKey,
	}
	nl.leveler = new(slog.LevelVar)
	nl.leveler.Set(s.leveler.Level())

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: nl.renameAttr,
		Level:       nl.leveler,
		AddSource:   true,
	}

	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	nl.slog = slog.New(handler).With(attrs...)

	s.mu.RUnlock()

	return nl
}

func (s *slogLogger) Init(opts ...logger.Option) error {
	s.mu.Lock()
	for _, o := range opts {
		o(&s.opts)
	}

	if v, ok := s.opts.Context.Value(levelKey{}).(string); ok && v != "" {
		s.levelKey = v
	}
	if v, ok := s.opts.Context.Value(messageKey{}).(string); ok && v != "" {
		s.messageKey = v
	}
	if v, ok := s.opts.Context.Value(sourceKey{}).(string); ok && v != "" {
		s.sourceKey = v
	}
	if v, ok := s.opts.Context.Value(timeKey{}).(string); ok && v != "" {
		s.timeKey = v
	}

	s.leveler = new(slog.LevelVar)
	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: s.renameAttr,
		Level:       s.leveler,
		AddSource:   true,
	}
	s.leveler.Set(loggerToSlogLevel(s.opts.Level))
	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	s.slog = slog.New(handler).With(s.opts.Fields...)

	slog.SetDefault(s.slog)

	s.mu.Unlock()

	return nil
}

func (s *slogLogger) Log(ctx context.Context, lvl logger.Level, attrs ...interface{}) {
	if !s.V(lvl) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), loggerToSlogLevel(lvl), fmt.Sprintf("%s", attrs[0]), pcs[0])
	// r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Logf(ctx context.Context, lvl logger.Level, msg string, attrs ...interface{}) {
	if !s.V(lvl) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), loggerToSlogLevel(lvl), fmt.Sprintf(msg, attrs...), pcs[0])
	// r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Info(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf("%s", attrs[0]), pcs[0])
	// r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Infof(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(msg, attrs...), pcs[0])
	// r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Debug(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.DebugLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf("%s", attrs[0]), pcs[0])
	// r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Debugf(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.DebugLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(msg, attrs...), pcs[0])
	// r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Trace(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.TraceLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelDebug-1, fmt.Sprintf("%s", attrs[0]), pcs[0])
	// r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Tracef(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.TraceLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelDebug-1, fmt.Sprintf(msg, attrs...), pcs[0])
	// r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Error(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.ErrorLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf("%s", attrs[0]), pcs[0])
	//	r.Add(attrs[1:]...)
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

func (s *slogLogger) Errorf(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.ErrorLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(msg, attrs...), pcs[0])
	// r.Add(attrs...)
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

func (s *slogLogger) Fatal(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.FatalLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError+1, fmt.Sprintf("%s", attrs[0]), pcs[0])
	// r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
	os.Exit(1)
}

func (s *slogLogger) Fatalf(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.FatalLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError+1, fmt.Sprintf(msg, attrs...), pcs[0])
	// r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
	os.Exit(1)
}

func (s *slogLogger) Warn(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.WarnLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf("%s", attrs[0]), pcs[0])
	// r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Warnf(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.WarnLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf(msg, attrs...), pcs[0])
	// r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Name() string {
	return s.opts.Name
}

func (s *slogLogger) String() string {
	return "slog"
}

func NewLogger(opts ...logger.Option) logger.Logger {
	s := &slogLogger{
		opts:       logger.NewOptions(opts...),
		sourceKey:  DefaultSourceKey,
		timeKey:    DefaultTimeKey,
		messageKey: DefaultMessageKey,
		levelKey:   DefaultLevelKey,
	}
	if v, ok := s.opts.Context.Value(levelKey{}).(string); ok && v != "" {
		s.levelKey = v
	}
	if v, ok := s.opts.Context.Value(messageKey{}).(string); ok && v != "" {
		s.messageKey = v
	}
	if v, ok := s.opts.Context.Value(sourceKey{}).(string); ok && v != "" {
		s.sourceKey = v
	}
	if v, ok := s.opts.Context.Value(timeKey{}).(string); ok && v != "" {
		s.timeKey = v
	}
	return s
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
