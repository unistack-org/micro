package logger

import (
	"log"
)

type stdLogger struct {
	l     Logger
	level Level
}

func NewStdLogger(l Logger, level Level) *log.Logger {
	return log.New(&stdLogger{l: l, level: level}, "" /* prefix */, 0 /* flags */)
}

func (sl *stdLogger) Write(p []byte) (int, error) {
	sl.l.Log(sl.l.Options().Context, sl.level, string(p))
	return len(p), nil
}
