package slog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"

	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/tracer"
)

var reTrace = regexp.MustCompile(`.*/slog/logger\.go.*\n`)

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
	mu      sync.RWMutex
}

func (s *slogLogger) Clone(opts ...logger.Option) logger.Logger {
	s.mu.RLock()
	options := s.opts

	for _, o := range opts {
		o(&options)
	}

	l := &slogLogger{
		opts: options,
	}

	l.leveler = new(slog.LevelVar)
	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: l.renameAttr,
		Level:       l.leveler,
		AddSource:   l.opts.AddSource,
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
	l := &slogLogger{opts: s.opts}
	l.leveler = new(slog.LevelVar)
	l.leveler.Set(s.leveler.Level())

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: l.renameAttr,
		Level:       l.leveler,
		AddSource:   l.opts.AddSource,
	}

	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	l.slog = slog.New(handler).With(attrs...)

	s.mu.RUnlock()

	return l
}

func (s *slogLogger) Init(opts ...logger.Option) error {
	s.mu.Lock()

	if len(s.opts.ContextAttrFuncs) == 0 {
		s.opts.ContextAttrFuncs = logger.DefaultContextAttrFuncs
	}

	for _, o := range opts {
		o(&s.opts)
	}

	s.leveler = new(slog.LevelVar)
	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: s.renameAttr,
		Level:       s.leveler,
		AddSource:   s.opts.AddSource,
	}
	s.leveler.Set(loggerToSlogLevel(s.opts.Level))
	handler := slog.NewJSONHandler(s.opts.Out, handleOpt)
	s.slog = slog.New(handler).With(s.opts.Fields...)

	s.mu.Unlock()

	return nil
}

func (s *slogLogger) Log(ctx context.Context, lvl logger.Level, attrs ...interface{}) {
	if !s.V(lvl) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), loggerToSlogLevel(lvl), fmt.Sprintf("%s", attrs[0]), pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	if s.opts.AddStacktrace && lvl == logger.ErrorLevel {
		stackInfo := make([]byte, 1024*1024)
		if stackSize := runtime.Stack(stackInfo, false); stackSize > 0 {
			traceLines := reTrace.Split(string(stackInfo[:stackSize]), -1)
			if len(traceLines) != 0 {
				attrs = append(attrs, slog.String(s.opts.StacktraceKey, traceLines[len(traceLines)-1]))
			}
		}
	}
	r.Add(attrs[1:]...)
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == s.opts.ErrorKey {
			if span, ok := tracer.SpanFromContext(ctx); ok {
				span.SetStatus(tracer.SpanStatusError, a.Value.String())
				return false
			}
		}
		return true
	})
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Logf(ctx context.Context, lvl logger.Level, msg string, attrs ...interface{}) {
	if !s.V(lvl) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), loggerToSlogLevel(lvl), msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	if s.opts.AddStacktrace && lvl == logger.ErrorLevel {
		stackInfo := make([]byte, 1024*1024)
		if stackSize := runtime.Stack(stackInfo, false); stackSize > 0 {
			traceLines := reTrace.Split(string(stackInfo[:stackSize]), -1)
			if len(traceLines) != 0 {
				attrs = append(attrs, (slog.String(s.opts.StacktraceKey, traceLines[len(traceLines)-1])))
			}
		}
	}
	r.Add(attrs[1:]...)
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == s.opts.ErrorKey {
			if span, ok := tracer.SpanFromContext(ctx); ok {
				span.SetStatus(tracer.SpanStatusError, a.Value.String())
				return false
			}
		}
		return true
	})
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Info(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelInfo, fmt.Sprintf("%s", attrs[0]), pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Infof(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.InfoLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelInfo, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Debug(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.DebugLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelDebug, fmt.Sprintf("%s", attrs[0]), pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Debugf(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.DebugLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelDebug, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Trace(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.TraceLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelDebug-1, fmt.Sprintf("%s", attrs[0]), pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Tracef(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.TraceLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelDebug-1, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Error(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.ErrorLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelError, fmt.Sprintf("%s", attrs[0]), pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	if s.opts.AddStacktrace {
		stackInfo := make([]byte, 1024*1024)
		if stackSize := runtime.Stack(stackInfo, false); stackSize > 0 {
			traceLines := reTrace.Split(string(stackInfo[:stackSize]), -1)
			if len(traceLines) != 0 {
				attrs = append(attrs, slog.String("stacktrace", traceLines[len(traceLines)-1]))
			}
		}
	}
	r.Add(attrs[1:]...)
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == s.opts.ErrorKey {
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
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelError, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	if s.opts.AddStacktrace {
		stackInfo := make([]byte, 1024*1024)
		if stackSize := runtime.Stack(stackInfo, false); stackSize > 0 {
			traceLines := reTrace.Split(string(stackInfo[:stackSize]), -1)
			if len(traceLines) != 0 {
				attrs = append(attrs, slog.String("stacktrace", traceLines[len(traceLines)-1]))
			}
		}
	}
	r.Add(attrs...)
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == s.opts.ErrorKey {
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
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelError+1, fmt.Sprintf("%s", attrs[0]), pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
	os.Exit(1)
}

func (s *slogLogger) Fatalf(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.FatalLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelError+1, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs...)
	_ = s.slog.Handler().Handle(ctx, r)
	os.Exit(1)
}

func (s *slogLogger) Warn(ctx context.Context, attrs ...interface{}) {
	if !s.V(logger.WarnLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelWarn, fmt.Sprintf("%s", attrs[0]), pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs[1:]...)
	_ = s.slog.Handler().Handle(ctx, r)
}

func (s *slogLogger) Warnf(ctx context.Context, msg string, attrs ...interface{}) {
	if !s.V(logger.WarnLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(s.opts.TimeFunc(), slog.LevelWarn, msg, pcs[0])
	for _, fn := range s.opts.ContextAttrFuncs {
		attrs = append(attrs, fn(ctx)...)
	}

	for idx, attr := range attrs {
		if ve, ok := attr.(error); ok && ve != nil {
			attrs[idx] = slog.String(s.opts.ErrorKey, ve.Error())
			break
		}
	}
	r.Add(attrs[1:]...)
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
		opts: logger.NewOptions(opts...),
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
