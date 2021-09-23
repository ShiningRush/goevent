# goevent
goevent is a simple pub/sub model in memory.

## usage
```go
package main
import (
    "github.com/shiningrush/goevent"
    "time"
)

type TestEvent struct {
	Key  string
	Key2 string
}

func (e *TestEvent) GetKey() string {
	return "test"
}

type TestEventHandler struct {
}

func (h *TestEventHandler) GetKey() string {
	return "test"
}

func (h *TestEventHandler) Handle(event goevent.Event) error {
	time.Sleep(time.Millisecond * 100)
	e := event.(*TestEvent)
	e.Key = "Handled"
	return nil
}

func main() {
    // sub
    err := goevent.Subscribe(&TestEventHandler{})
    if err != nil {
        panic(err)
    }

    e := &TestEvent{}
    // pub
    err = goevent.Publish(e)
    if err != nil {
        panic(err)
    }
    // pub async
    goevent.PublishAsync(e)
}
```