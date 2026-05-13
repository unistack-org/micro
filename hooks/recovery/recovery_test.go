package recovery

import (
	"context"
	"testing"

	"go.unistack.org/micro/v5/errors"
	"go.unistack.org/micro/v5/server"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	if opts.ServerHandlerFn == nil {
		t.Error("expected non-nil ServerHandlerFn")
	}
}

func TestNewHook(t *testing.T) {
	h := NewHook()
	if h == nil {
		t.Error("expected non-nil hook")
	}
}

func TestNewHookWithOption(t *testing.T) {
	custom := func(ctx context.Context, req server.Request, rsp any, err error) error {
		return errors.BadRequest("custom", "error: %v", err)
	}
	h := NewHook(ServerHandlerFunc(custom))
	if h == nil {
		t.Error("expected non-nil hook")
	}
	_ = custom
}

func TestHookServerHandlerPanic(t *testing.T) {
	h := NewHook()
	next := func(ctx context.Context, req server.Request, rsp any) error {
		panic("intentional panic")
	}
	wrapped := h.ServerHandler(next)
	err := wrapped(context.Background(), nil, nil)
	if err == nil {
		t.Error("expected error from panic")
	}
}

func TestHookServerHandlerNoPanic(t *testing.T) {
	h := NewHook()
	next := func(ctx context.Context, req server.Request, rsp any) error {
		return nil
	}
	wrapped := h.ServerHandler(next)
	err := wrapped(context.Background(), nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDefaultServerHandlerFn(t *testing.T) {
	err := DefaultServerHandlerFn(context.Background(), nil, nil, nil)
	if err == nil {
		t.Error("expected error")
	}
}