package micro

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/fsm"
)

func TestAs(t *testing.T) {
	var b *bro
	broTarget := &bro{name: "kafka"}
	fsmTarget := &fsmT{name: "fsm"}

	testCases := []struct {
		b      any
		target any
		match  bool
		want   any
	}{
		{
			broTarget,
			&b,
			true,
			broTarget,
		},
		{
			nil,
			&b,
			false,
			nil,
		},
		{
			fsmTarget,
			&b,
			false,
			nil,
		},
	}
	for i, tc := range testCases {
		name := fmt.Sprintf("%d:As(Errorf(..., %v), %v)", i, tc.b, tc.target)
		// Clear the target pointer, in case it was set in a previous test.
		rtarget := reflect.ValueOf(tc.target)
		rtarget.Elem().Set(reflect.Zero(reflect.TypeOf(tc.target).Elem()))
		t.Run(name, func(t *testing.T) {
			match := As(tc.b, tc.target)
			if match != tc.match {
				t.Fatalf("match: got %v; want %v", match, tc.match)
			}
			if !match {
				return
			}
			if got := rtarget.Elem().Interface(); got != tc.want {
				t.Fatalf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

type bro struct {
	name string
}

func (p *bro) Name() string { return p.name }

func (p *bro) Live() bool { return true }

func (p *bro) Ready() bool { return true }

func (p *bro) Health() bool { return true }

func (p *bro) Init(opts ...broker.Option) error { return nil }

// Options returns broker options
func (p *bro) Options() broker.Options { return broker.Options{} }

// Address return configured address
func (p *bro) Address() string { return "" }

// Connect connects to broker
func (p *bro) Connect(ctx context.Context) error { return nil }

// Disconnect disconnect from broker
func (p *bro) Disconnect(ctx context.Context) error { return nil }

// Publish message, msg can be single broker.Message or []broker.Message
func (p *bro) Publish(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	return nil
}

// BatchPublish messages to broker with multiple topics
func (p *bro) BatchPublish(ctx context.Context, msgs []*broker.Message, opts ...broker.PublishOption) error {
	return nil
}

// BatchSubscribe subscribes to topic messages via handler
func (p *bro) BatchSubscribe(ctx context.Context, topic string, h broker.BatchHandler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return nil, nil
}

// Subscribe subscribes to topic message via handler
func (p *bro) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return nil, nil
}

// String type of broker
func (p *bro) String() string { return p.name }

type fsmT struct {
	name string
}

func (f *fsmT) Start(ctx context.Context, a interface{}, o ...Option) (interface{}, error) {
	return nil, nil
}
func (f *fsmT) Current() string                  { return f.name }
func (f *fsmT) Reset()                           {}
func (f *fsmT) State(s string, sf fsm.StateFunc) {}
