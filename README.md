# goevent
goevent is a simple pub/sub model in memory.

## usage
```go


package main

import (
	"context"
	"fmt"
	"github.com/shiningrush/goevent"
)

type DemoEventA struct {
	Key string
}

func (e *DemoEventA) Topic() []string {
	return []string{"topicX", "topicY"}
}

type DemoEventB struct {
	Key string
}

func (e *DemoEventB) Topic() []string {
	return []string{"topicY", "topicZ"}
}

type TestEventHandler struct {
}

func (h *TestEventHandler) Topic() []string {
	return []string{"topicX", "topicY", "topicZ"}
}

func (h *TestEventHandler) Handle(cxt context.Context, event goevent.Event) {
	if eventA, ok := event.(*DemoEventA); ok {
		fmt.Println("get event a, key: ", eventA.Key)
	}
	if eventB, ok := event.(*DemoEventB); ok {
		fmt.Println("get event b, key: ", eventB.Key)
	}
}

func main() {
	// sub
	err := goevent.Subscribe(&TestEventHandler{})
	if err != nil {
		panic(err)
	}

	var e goevent.Event
	e = &DemoEventA{Key: "value a"}
	// should print "get event a, key:  value a" twice, because handler subscribe all topics
	goevent.PublishSync(context.TODO(), e)

	e = &DemoEventB{Key: "value b"}
	goevent.PublishSync(context.TODO(), e)
}
```
