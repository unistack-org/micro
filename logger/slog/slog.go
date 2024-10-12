package slog

import (
	"context"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"

	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/semconv"
	"go.unistack.org/micro/v3/tracer"
)

const (
	badKey = "!BADKEY"
	// defaultCallerSkipCount used by logger
	defaultCallerSkipCount = 3
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
	leveler *slog.LevelVar
	handler slog.Handler
	opts    logger.Options
	mu      sync.RWMutex
}

func (s *slogLogger) Clone(opts ...logger.Option) logger.Logger {
	s.mu.RLock()
	options := s.opts
	s.mu.RUnlock()

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
	l.handler = slog.New(slog.NewJSONHandler(options.Out, handleOpt)).With(options.Fields...).Handler()

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
	level := s.leveler.Level()
	options := s.opts
	s.mu.RUnlock()

	l := &slogLogger{opts: options}
	l.leveler = new(slog.LevelVar)
	l.leveler.Set(level)

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: l.renameAttr,
		Level:       l.leveler,
		AddSource:   l.opts.AddSource,
	}

	l.handler = slog.New(slog.NewJSONHandler(l.opts.Out, handleOpt)).With(attrs...).Handler()

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
	s.handler = slog.New(slog.NewJSONHandler(s.opts.Out, handleOpt)).With(s.opts.Fields...).Handler()
	s.mu.Unlock()

	return nil
}

func (s *slogLogger) Log(ctx context.Context, lvl logger.Level, msg string, attrs ...interface{}) {
	s.printLog(ctx, lvl, msg, attrs...)
}

func (s *slogLogger) Info(ctx context.Context, msg string, attrs ...interface{}) {
	s.printLog(ctx, logger.InfoLevel, msg, attrs...)
}

func (s *slogLogger) Debug(ctx context.Context, msg string, attrs ...interface{}) {
	s.printLog(ctx, logger.DebugLevel, msg, attrs...)
}

func (s *slogLogger) Trace(ctx context.Context, msg string, attrs ...interface{}) {
	s.printLog(ctx, logger.TraceLevel, msg, attrs...)
}

func (s *slogLogger) Error(ctx context.Context, msg string, attrs ...interface{}) {
	s.printLog(ctx, logger.ErrorLevel, msg, attrs...)
}

func (s *slogLogger) Fatal(ctx context.Context, msg string, attrs ...interface{}) {
	s.printLog(ctx, logger.FatalLevel, msg, attrs...)
	os.Exit(1)
}

func (s *slogLogger) Warn(ctx context.Context, msg string, attrs ...interface{}) {
	s.printLog(ctx, logger.WarnLevel, msg, attrs...)
}

func (s *slogLogger) Name() string {
	return s.opts.Name
}

func (s *slogLogger) String() string {
	return "slog"
}

func (s *slogLogger) printLog(ctx context.Context, lvl logger.Level, msg string, attrs ...interface{}) {
	if !s.V(lvl) {
		return
	}

	s.opts.Meter.Counter(semconv.LoggerMessageTotal, "level", lvl.String()).Inc()

	attrs = prepareAttributes(attrs)

	for _, fn := range s.opts.ContextAttrFuncs {
		a := prepareAttributes(fn(ctx))
		attrs = append(attrs, a...)
	}

	for _, attr := range attrs {
		if ve, hasErr := attr.(error); hasErr && ve != nil {
			attrs = append(attrs, slog.String(s.opts.ErrorKey, ve.Error()))
			if span, ok := tracer.SpanFromContext(ctx); ok {
				span.SetStatus(tracer.SpanStatusError, ve.Error())
			}
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

	var pcs [1]uintptr
	runtime.Callers(s.opts.CallerSkipCount, pcs[:]) // skip [Callers, printLog, LogLvlMethod]
	r := slog.NewRecord(s.opts.TimeFunc(), loggerToSlogLevel(lvl), msg, pcs[0])
	r.Add(attrs...)
	_ = s.handler.Handle(ctx, r)
}

func NewLogger(opts ...logger.Option) logger.Logger {
	s := &slogLogger{
		opts: logger.NewOptions(opts...),
	}
	s.opts.CallerSkipCount = defaultCallerSkipCount

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

func prepareAttributes(attrs []interface{}) []interface{} {
	if len(attrs)%2 == 1 {
		attrs = append(attrs, badKey)
		attrs[len(attrs)-1], attrs[len(attrs)-2] = attrs[len(attrs)-2], attrs[len(attrs)-1]
	}

	return attrs
}
