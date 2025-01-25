package fsm

import (
	"context"
	"fmt"
	"testing"

	"go.unistack.org/micro/v4/logger"
)

func TestFSMStart(t *testing.T) {
	ctx := context.TODO()

	if err := logger.DefaultLogger.Init(); err != nil {
		t.Fatal(err)
	}

	wrapper := func(next StateFunc) StateFunc {
		return func(sctx context.Context, s State, opts ...StateOption) (State, error) {
			sctx = logger.NewContext(sctx, logger.DefaultLogger.Fields("state", s.Name()))
			return next(sctx, s, opts...)
		}
	}

	f := NewFSM(InitialState("1"), WrapState(wrapper))
	f1 := func(sctx context.Context, s State, opts ...StateOption) (State, error) {
		_, ok := logger.FromContext(sctx)
		if !ok {
			t.Fatal("f1 context does not have logger")
		}
		args := s.Body().(map[string]interface{})
		if v, ok := args["request"].(string); !ok || v == "" {
			return nil, fmt.Errorf("empty request")
		}
		return &state{name: "", body: map[string]interface{}{"response": "state1"}}, nil
	}
	f2 := func(sctx context.Context, s State, opts ...StateOption) (State, error) {
		_, ok := logger.FromContext(sctx)
		if !ok {
			t.Fatal("f2 context does not have logger")
		}
		args := s.Body().(map[string]interface{})
		if v, ok := args["response"].(string); !ok || v == "" {
			return nil, fmt.Errorf("empty response")
		}
		return &state{name: "", body: map[string]interface{}{"response": "state2"}}, nil
	}
	f3 := func(sctx context.Context, s State, opts ...StateOption) (State, error) {
		_, ok := logger.FromContext(sctx)
		if !ok {
			t.Fatal("f3 context does not have logger")
		}
		args := s.Body().(map[string]interface{})
		if v, ok := args["response"].(string); !ok || v == "" {
			return nil, fmt.Errorf("empty response")
		}
		return &state{name: StateEnd, body: map[string]interface{}{"response": "state3"}}, nil
	}
	f.State("1", f1)
	f.State("2", f2)
	f.State("3", f3)
	rsp, err := f.Start(ctx, map[string]interface{}{"request": "state"})
	if err != nil {
		t.Fatal(err)
	}
	args := rsp.(map[string]interface{})
	if v, ok := args["response"].(string); !ok || v == "" {
		t.Fatalf("nil rsp: %#+v", args)
	} else if v != "state3" {
		t.Fatalf("invalid rsp %#+v", args)
	}
}
