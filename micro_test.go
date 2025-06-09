package micro

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/fsm"
	"go.unistack.org/micro/v4/metadata"
)

func TestAs(t *testing.T) {
	var b *bro
	broTarget := &bro{name: "kafka"}
	fsmTarget := &fsmT{name: "fsm"}

	testCases := []struct {
		b      any
		target any
		want   any

		match bool
	}{
		{
			b:      broTarget,
			target: &b,
			match:  true,
			want:   broTarget,
		},
		{
			b:      nil,
			target: &b,
			match:  false,
			want:   nil,
		},
		{
			b:      fsmTarget,
			target: &b,
			match:  false,
			want:   nil,
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

var _ broker.Broker = (*bro)(nil)

type bro struct {
	name string
}

func (p *bro) Name() string { return p.name }

func (p *bro) Live() bool { return true }

func (p *bro) Ready() bool { return true }

func (p *bro) Health() bool { return true }

func (p *bro) Init(_ ...broker.Option) error { return nil }

// Options returns broker options
func (p *bro) Options() broker.Options { return broker.Options{} }

// Address return configured address
func (p *bro) Address() string { return "" }

// Connect connects to broker
func (p *bro) Connect(_ context.Context) error { return nil }

// Disconnect disconnect from broker
func (p *bro) Disconnect(_ context.Context) error { return nil }

// NewMessage creates new message
func (p *bro) NewMessage(_ context.Context, _ metadata.Metadata, _ interface{}, _ ...broker.MessageOption) (broker.Message, error) {
	return nil, nil
}

// Publish message, msg can be single broker.Message or []broker.Message
func (p *bro) Publish(_ context.Context, _ string, _ ...broker.Message) error {
	return nil
}

// Subscribe subscribes to topic message via handler
func (p *bro) Subscribe(_ context.Context, _ string, _ interface{}, _ ...broker.SubscribeOption) (broker.Subscriber, error) {
	return nil, nil
}

// String type of broker
func (p *bro) String() string { return p.name }

type fsmT struct {
	name string
}

func (f *fsmT) Start(_ context.Context, _ interface{}, _ ...Option) (interface{}, error) {
	return nil, nil
}
func (f *fsmT) Current() string                 { return f.name }
func (f *fsmT) Reset()                          {}
func (f *fsmT) State(_ string, _ fsm.StateFunc) {}
