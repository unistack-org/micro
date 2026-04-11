package slog

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"

	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/semconv"
	"go.unistack.org/micro/v4/tracer"
	pool "go.unistack.org/micro/v4/util/xpool"
)

const (
	badKey = "!BADKEY"
	// defaultCallerSkipCount used by logger
	defaultCallerSkipCount = 3
	timeFormat             = "2006-01-02T15:04:05.000000000Z07:00"
)

var reTrace = regexp.MustCompile(`.*/slog/logger\.go.*\n`)

// pool for slog.Attr slice to reduce allocations
var recordPool = pool.NewPool[[]slog.Attr](func() []slog.Attr {
	return make([]slog.Attr, 0, 32)
}, 1024)

// pool for stack trace buffers to reduce allocations
var stackPool = pool.NewBytePool(1024 * 1024)

var (
	traceValue = slog.StringValue("trace")
	debugValue = slog.StringValue("debug")
	infoValue  = slog.StringValue("info")
	warnValue  = slog.StringValue("warn")
	errorValue = slog.StringValue("error")
	fatalValue = slog.StringValue("fatal")
	noneValue  = slog.StringValue("none")
)

type wrapper struct {
	h     slog.Handler
	level int64
}

func (h *wrapper) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= slog.Level(atomic.LoadInt64(&h.level))
}

func (h *wrapper) Handle(ctx context.Context, rec slog.Record) error {
	return h.h.Handle(ctx, rec)
}

func (h *wrapper) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &wrapper{
		h:     h.h.WithAttrs(attrs),
		level: atomic.LoadInt64(&h.level),
	}
}

func (h *wrapper) WithGroup(name string) slog.Handler {
	return &wrapper{
		h:     h.h.WithGroup(name),
		level: atomic.LoadInt64(&h.level),
	}
}

func (s *slogLogger) renameAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.SourceKey:
		source := a.Value.Any().(*slog.Source)
		a.Value = slog.StringValue(source.File + ":" + strconv.Itoa(source.Line))
		a.Key = s.opts.SourceKey
	case slog.TimeKey:
		a.Key = s.opts.TimeKey
		a.Value = slog.StringValue(a.Value.Time().Format(timeFormat))
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
		case lvl >= logger.NoneLevel:
			a.Value = noneValue
		default:
			a.Value = infoValue
		}
	}

	return a
}

type slogLogger struct {
	handler *wrapper
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

	if len(options.ContextAttrFuncs) == 0 {
		options.ContextAttrFuncs = logger.DefaultContextAttrFuncs
	}

	attrs := recordPool.Get()
	attrs = attrs[:0]
	defer recordPool.Put(attrs)

	s.argsAttrs(attrs, options.Fields...)
	l := &slogLogger{
		handler: &wrapper{
			h:     s.handler.h.WithAttrs(attrs),
			level: atomic.LoadInt64(&s.handler.level),
		},
		opts: options,
	}
	atomic.StoreInt64(&l.handler.level, int64(loggerToSlogLevel(options.Level)))

	return l
}

func (s *slogLogger) V(level logger.Level) bool {
	s.mu.Lock()
	v := s.opts.Level.Enabled(level)
	s.mu.Unlock()
	return v
}

func (s *slogLogger) Level(level logger.Level) {
	atomic.StoreInt64(&s.handler.level, int64(loggerToSlogLevel(level)))
	s.mu.Lock()
	s.opts.Level = level
	s.mu.Unlock()
}

func (s *slogLogger) Options() logger.Options {
	return s.opts
}

func (s *slogLogger) Fields(fields ...interface{}) logger.Logger {
	s.mu.RLock()
	options := s.opts
	s.mu.RUnlock()

	l := &slogLogger{opts: options}
	logger.WithAddFields(fields...)(&l.opts)

	if len(options.ContextAttrFuncs) == 0 {
		options.ContextAttrFuncs = logger.DefaultContextAttrFuncs
	}

	attrs := recordPool.Get()
	attrs = attrs[:0]
	defer recordPool.Put(attrs)

	s.argsAttrs(attrs, fields...)
	l.handler = &wrapper{
		h:     s.handler.h.WithAttrs(attrs),
		level: atomic.LoadInt64(&s.handler.level),
	}
	atomic.StoreInt64(&l.handler.level, int64(loggerToSlogLevel(l.opts.Level)))

	return l
}

func (s *slogLogger) Init(opts ...logger.Option) error {
	s.mu.Lock()

	for _, o := range opts {
		o(&s.opts)
	}

	if len(s.opts.ContextAttrFuncs) == 0 {
		s.opts.ContextAttrFuncs = logger.DefaultContextAttrFuncs
	}

	handleOpt := &slog.HandlerOptions{
		ReplaceAttr: s.renameAttr,
		Level:       loggerToSlogLevel(logger.TraceLevel),
		AddSource:   s.opts.AddSource,
	}

	attrs := recordPool.Get()
	attrs = attrs[:0]
	defer recordPool.Put(attrs)

	s.argsAttrs(attrs, s.opts.Fields...)

	var h slog.Handler
	if s.opts.Context != nil {
		if v, ok := s.opts.Context.Value(handlerKey{}).(slog.Handler); ok && v != nil {
			h = v
		}

		if fn := s.opts.Context.Value(handlerFnKey{}); fn != nil {
			if rfn := reflect.ValueOf(fn); rfn.Kind() == reflect.Func {
				if ret := rfn.Call([]reflect.Value{reflect.ValueOf(s.opts.Out), reflect.ValueOf(handleOpt)}); len(ret) == 1 {
					if iface, ok := ret[0].Interface().(slog.Handler); ok && iface != nil {
						h = iface
					}
				}
			}
		}
	}

	if h == nil {
		h = slog.NewJSONHandler(s.opts.Out, handleOpt)
	}

	s.handler = &wrapper{
		h: h.WithAttrs(attrs),
	}

	atomic.StoreInt64(&s.handler.level, int64(loggerToSlogLevel(s.opts.Level)))
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
	if closer, ok := s.opts.Out.(io.Closer); ok {
		_ = closer.Close()
	}
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

func (s *slogLogger) printLog(ctx context.Context, lvl logger.Level, msg string, args ...interface{}) {
	if !s.V(lvl) {
		return
	}
	var argError error

	s.opts.Meter.Counter(semconv.LoggerMessageTotal, "level", lvl.String()).Inc()

	attrs := recordPool.Get()
	attrs = attrs[:0]
	defer recordPool.Put(attrs)

	var err error
	attrs, err = s.argsAttrs(attrs, args)
	if err != nil {
		argError = err
	}
	if argError != nil {
		if span, ok := tracer.SpanFromContext(ctx); ok {
			span.SetStatus(tracer.SpanStatusError, argError.Error())
		}
	}

	for _, fn := range s.opts.ContextAttrFuncs {
		var ctxErr error
		attrs, ctxErr = s.argsAttrs(attrs, fn(ctx))
		if ctxErr != nil {
			argError = ctxErr
		}
	}
	if argError != nil {
		if span, ok := tracer.SpanFromContext(ctx); ok {
			span.SetStatus(tracer.SpanStatusError, argError.Error())
		}
	}

	if s.opts.AddStacktrace && (lvl == logger.FatalLevel || lvl == logger.ErrorLevel) {
		stackBuf := stackPool.Get()
		*stackBuf = (*stackBuf)[:0]
		defer stackPool.Put(stackBuf)
		if stackSize := runtime.Stack(*stackBuf, false); stackSize > 0 {
			traceLines := reTrace.Split(string((*stackBuf)[:stackSize]), -1)
			if len(traceLines) != 0 {
				attrs = append(attrs, slog.String(s.opts.StacktraceKey, traceLines[len(traceLines)-1]))
			}
		}
	}

	var pcs uintptr

	if s.opts.AddCaller {
		var caller [1]uintptr
		runtime.Callers(s.opts.CallerSkipCount, caller[:]) // skip [Callers, printLog, LogLvlMethod]
		pcs = caller[0]
	}

	r := slog.NewRecord(s.opts.TimeFunc(), loggerToSlogLevel(lvl), msg, pcs)
	r.AddAttrs(attrs...)

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
	case logger.NoneLevel:
		return slog.LevelError + 2
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
	case slog.LevelError + 2:
		return logger.NoneLevel
	default:
		return logger.InfoLevel
	}
}

func (s *slogLogger) argsAttrs(attrs []slog.Attr, args ...interface{}) ([]slog.Attr, error) {
	var err error

	for idx := 0; idx < len(args); idx++ {
		switch arg := args[idx].(type) {
		case slog.Attr:
			attrs = append(attrs, arg)
		case string:
			if idx+1 < len(args) {
				attrs = append(attrs, slog.Any(arg, args[idx+1]))
				idx++
			} else {
				attrs = append(attrs, slog.String(badKey, arg))
			}
		case error:
			attrs = append(attrs, slog.String(s.opts.ErrorKey, arg.Error()))
			err = arg
		}
	}

	return attrs, err
}

type handlerKey struct{}

func WithHandler(h slog.Handler) logger.Option {
	return logger.SetOption(handlerKey{}, h)
}

type handlerFnKey struct{}

func WithHandlerFunc(fn any) logger.Option {
	return logger.SetOption(handlerFnKey{}, fn)
}
