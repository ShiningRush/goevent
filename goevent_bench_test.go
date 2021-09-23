package goevent

import (
	"context"
	"testing"
)

func BenchmarkInMemoryEventBus_Close(b *testing.B) {
	testBus := NewInMemoryEventBus()

	_ = testBus.Subscribe(&TestEventHandler{})
	_ = testBus.Subscribe(&SecTestEventHandler{})
	e := &TestEvent{}
	for i := 0; i < b.N; i++ {
		testBus.Publish(e)
	}
	testBus.Close()
}

func BenchmarkInMemoryEventBus(b *testing.B) {
	_ = Subscribe(&TestEventHandler{})
	_ = Subscribe(&SecTestEventHandler{})
	e := &TestEvent{}
	for i := 0; i < b.N; i++ {
		Publish(e)
	}
}

func BenchmarkInMemoryEventBus_Sync(b *testing.B) {
	_ = Subscribe(&TestEventHandler{})
	_ = Subscribe(&SecTestEventHandler{})

	e := &TestEvent{}
	for i := 0; i < b.N; i++ {
		PublishSync(context.TODO(), e)
	}
}

func BenchmarkDirectCall(b *testing.B) {
	e := &TestEvent{}
	h1 := TestEventHandler{}
	h2 := SecTestEventHandler{}
	for i := 0; i < b.N; i++ {
		h1.Handle(context.TODO(), e)
		h2.Handle(context.TODO(), e)
	}
}
