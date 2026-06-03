# goevent (Deprecated)

goevent is a simple in-memory pub/sub library.

This project is deprecated. Please use `github.com/shiningrush/goext/runx/eventx` in new code.

The `github.com/shiningrush/goevent` package is kept only as a compatibility layer and now forwards to the latest `goext/runx/eventx` implementation.

## migration

Recommended import:

```go
import "github.com/shiningrush/goext/runx/eventx"
```

Old code:

```go
import "github.com/shiningrush/goevent"
```

New code:

```go
import "github.com/shiningrush/goext/runx/eventx"
```

Most public APIs map directly:

- `goevent.Event` -> `eventx.Event`
- `goevent.EventHandler` -> `eventx.EventHandler`
- `goevent.EventBus` -> `eventx.EventBus`
- `goevent.NewInMemoryEventBus` -> `eventx.NewInMemoryEventBus`
- `goevent.Subscribe` -> `eventx.Subscribe`
- `goevent.Publish` -> `eventx.Publish`
- `goevent.PublishSync` -> `eventx.PublishSync`
- `goevent.Close` -> `eventx.Close`

## usage
```go


package main

import (
	"context"
	"fmt"

	"github.com/shiningrush/goext/runx/eventx"
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

func (h *TestEventHandler) Handle(cxt context.Context, event eventx.Event) {
	if eventA, ok := event.(*DemoEventA); ok {
		fmt.Println("get event a, key: ", eventA.Key)
	}
	if eventB, ok := event.(*DemoEventB); ok {
		fmt.Println("get event b, key: ", eventB.Key)
	}
}

func main() {
	// sub
	err := eventx.Subscribe(&TestEventHandler{})
	if err != nil {
		panic(err)
	}

	var e eventx.Event
	e = &DemoEventA{Key: "value a"}
	// should print "get event a, key:  value a" twice, because handler subscribe all topics
	eventx.PublishSync(context.TODO(), e)

	e = &DemoEventB{Key: "value b"}
	eventx.PublishSync(context.TODO(), e)
}
```
