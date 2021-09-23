package goevent

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestEvent struct {
	Key  string
	Key2 string
}

func (e *TestEvent) Topic() []string {
	return []string{"test1", "test2"}
}

type TestEventHandler struct {
}

func (h *TestEventHandler) Topic() []string {
	return []string{"test1"}
}

func (h *TestEventHandler) Handle(cxt context.Context, event Event) {
	time.Sleep(time.Millisecond * 100)
	e := event.(*TestEvent)
	e.Key = "Handled"
}

type SecTestEventHandler struct {
}

func (h *SecTestEventHandler) Topic() []string {
	return []string{"test2"}
}

func (h *SecTestEventHandler) Handle(cxt context.Context, event Event) {
	time.Sleep(time.Millisecond * 100)
	e := event.(*TestEvent)
	e.Key2 = "Handled2"
}

func TestPublish(t *testing.T) {
	err := Subscribe(&TestEventHandler{})
	err = Subscribe(&SecTestEventHandler{})

	assert.NoError(t, err)
	e := &TestEvent{Key: "UnHandle", Key2: "UnHandle2"}
	Publish(e)
	assert.NoError(t, err)
	assert.Equal(t, "UnHandle", e.Key)
	assert.Equal(t, "UnHandle2", e.Key2)
	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, "Handled", e.Key)
	assert.Equal(t, "Handled2", e.Key2)
}

func TestPublishSync(t *testing.T) {
	err := Subscribe(&TestEventHandler{})
	err = Subscribe(&SecTestEventHandler{})

	assert.NoError(t, err)
	e := &TestEvent{Key: "UnHandle", Key2: "UnHandle2"}
	PublishSync(context.TODO(), e)
	assert.Equal(t, "Handled", e.Key)
	assert.Equal(t, "Handled2", e.Key2)
}

func TestGetEventKey(t *testing.T) {
	bus := NewInMemoryEventBus()

	key := bus.getEventTopic(reflect.TypeOf(bus))
	assert.Equal(t, reflect.TypeOf(bus), key)
}

func TestClose(t *testing.T) {
	testBus := NewInMemoryEventBus()
	err := testBus.Subscribe(&TestEventHandler{})
	err = testBus.Subscribe(&SecTestEventHandler{})

	assert.NoError(t, err)
	e := &TestEvent{Key: "UnHandle", Key2: "UnHandle2"}
	testBus.Publish(e)
	assert.NoError(t, err)
	assert.Equal(t, "UnHandle", e.Key)
	assert.Equal(t, "UnHandle2", e.Key2)
	testBus.Close()
	assert.Equal(t, "Handled", e.Key)
	assert.Equal(t, "Handled2", e.Key2)

	defer func() {
		err := recover()
		assert.Equal(t, "event bus is already closed", err)
	}()
	testBus.Publish(e)
}
