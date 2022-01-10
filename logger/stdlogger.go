package logger

import (
	"bytes"
	"log"
)

type stdLogger struct {
	l     Logger
	level Level
}

// NewStdLogger returns new *log.Logger baked by logger.Logger implementation
func NewStdLogger(l Logger, level Level) *log.Logger {
	return log.New(&stdLogger{l: l, level: level}, "" /* prefix */, 0 /* flags */)
}

func (sl *stdLogger) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	sl.l.Log(sl.l.Options().Context, sl.level, string(p))
	return len(p), nil
}

// RedirectStdLogger replace *log.Logger with logger.Logger implementation
func RedirectStdLogger(l Logger, level Level) func() {
	flags := log.Flags()
	prefix := log.Prefix()
	writer := log.Writer()
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(&stdLogger{l: l, level: level})
	return func() {
		log.SetFlags(flags)
		log.SetPrefix(prefix)
		log.SetOutput(writer)
	}
}
