package mock_test

import (
	"context"
	"fmt"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/broker/mock"
	"go.unistack.org/micro/v5/metadata"
)

// Example demonstrates the full lifecycle: connect, subscribe, inject a message,
// publish, unsubscribe, disconnect, and verify all expectations were met.
func Example() {
	ctx := context.Background()
	b := mock.NewMockBroker()

	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectPublish("results")
	b.ExpectUnsubscribe("orders")
	b.ExpectDisconnect()

	_ = b.Connect(ctx)

	sub, _ := b.Subscribe(ctx, "orders", func(m broker.Message) error {
		fmt.Printf("received: %s\n", m.Body())
		return nil
	})

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{}, []byte(`{"id":1}`), nil)
	_ = b.InjectMessage(ctx, "orders", msg)

	_ = b.Publish(ctx, "results")
	_ = sub.Unsubscribe(ctx)
	_ = b.Disconnect(ctx)

	if err := b.ExpectationsWereMet(); err != nil {
		fmt.Println("unmet expectations:", err)
	}
	// Output:
	// received: {"id":1}
}

// ExampleMockBroker_ExpectPublish demonstrates testing that application code
// publishes to the correct topic.
func ExampleMockBroker_ExpectPublish() {
	ctx := context.Background()
	b := mock.NewMockBroker()

	b.ExpectConnect()
	b.ExpectPublish("orders")
	b.ExpectDisconnect()

	_ = b.Connect(ctx)

	// application code under test
	if err := b.Publish(ctx, "orders"); err != nil {
		fmt.Println("publish error:", err)
	}

	_ = b.Disconnect(ctx)

	if err := b.ExpectationsWereMet(); err != nil {
		fmt.Println("unmet:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}

// ExampleMockBroker_InjectMessage demonstrates testing a subscriber handler
// by injecting a pre-built message.
func ExampleMockBroker_InjectMessage() {
	ctx := context.Background()
	b := mock.NewMockBroker()

	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectUnsubscribe("orders")
	b.ExpectDisconnect()

	_ = b.Connect(ctx)

	sub, _ := b.Subscribe(ctx, "orders", func(m broker.Message) error {
		fmt.Printf("handler got: %s\n", m.Body())
		return nil
	})

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{}, []byte(`{"id":7}`), nil)
	_ = b.InjectMessage(ctx, "orders", msg)

	_ = sub.Unsubscribe(ctx)
	_ = b.Disconnect(ctx)
	// Output:
	// handler got: {"id":7}
}
