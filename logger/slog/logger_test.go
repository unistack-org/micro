package slog

import (
	"bytes"
	"context"
	"testing"

	"go.unistack.org/micro/v5/logger"
)

func TestSlogLogger_Level(t *testing.T) {
	l := NewLogger()
	_ = l.Init()
	l.Level(logger.InfoLevel)
	if !l.V(logger.InfoLevel) {
		t.Error("expected V(InfoLevel) to be true")
	}
	if l.V(logger.DebugLevel) {
		t.Error("expected V(DebugLevel) to be false")
	}
}

func TestSlogLogger_Fields(t *testing.T) {
	l := NewLogger()
	_ = l.Init()
	l2 := l.Fields("key", "value")
	if l2 == nil {
		t.Error("expected non-nil logger from Fields")
	}
}

func TestSlogLogger_Clone(t *testing.T) {
	l := NewLogger()
	_ = l.Init()
	cloned := l.Clone()
	if cloned == nil {
		t.Error("expected non-nil cloned logger")
	}
}

func TestSlogLogger_Log(t *testing.T) {
	l := NewLogger()
	_ = l.Init()
	ctx := context.Background()
	l.Log(ctx, logger.InfoLevel, "test message", "key", "value")
}

func TestSlogLogger_Write(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(logger.WithOutput(&buf))
	_ = l.Init()
	ctx := context.Background()
	l.Info(ctx, "test message")
	if buf.Len() == 0 {
		t.Error("expected output to be written")
	}
}
