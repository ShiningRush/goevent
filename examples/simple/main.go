package main

import (
	"context"
	"fmt"
	"github.com/shiningrush/goevent"
	"time"
)

type TestEvent struct {
	Key string
}

func (e *TestEvent) Topic() []string {
	return []string{"test"}
}

type TestEventHandler struct {
}

func (h *TestEventHandler) Topic() []string {
	return []string{"test"}
}

func (h *TestEventHandler) Handle(cxt context.Context, event goevent.Event) {
	time.Sleep(time.Millisecond * 100)
	e := event.(*TestEvent)
	fmt.Println("event key: ", e.Key)
}

func main() {
	// sub
	err := goevent.Subscribe(&TestEventHandler{})
	if err != nil {
		panic(err)
	}

	e := &TestEvent{
		Key: time.Now().String(),
	}
	// publish async
	goevent.Publish(e)

	// publish sync
	goevent.PublishSync(context.TODO(), e)
}
