package broker

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go.unistack.org/micro/v4/options"
)

type bro struct {
	name string
}

func (p *bro) Name() string                      { return p.name }
func (p *bro) Init(opts ...options.Option) error { return nil }

// Options returns broker options
func (p *bro) Options() Options { return Options{} }

// Address return configured address
func (p *bro) Address() string { return "" }

// Connect connects to broker
func (p *bro) Connect(ctx context.Context) error { return nil }

// Disconnect disconnect from broker
func (p *bro) Disconnect(ctx context.Context) error { return nil }

// Publish message, msg can be single broker.Message or []broker.Message
func (p *bro) Publish(ctx context.Context, msg interface{}, opts ...options.Option) error { return nil }

// Subscribe subscribes to topic message via handler
func (p *bro) Subscribe(ctx context.Context, topic string, handler interface{}, opts ...options.Option) (Subscriber, error) {
	return nil, nil
}

// String type of broker
func (p *bro) String() string { return p.name }

func TestAs(t *testing.T) {
	var b *bro
	broTarget := &bro{name: "kafka"}

	testCases := []struct {
		b      Broker
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
