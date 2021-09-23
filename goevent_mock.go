package goevent

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// MockEventBus used to execute unit test
type MockEventBus struct {
	mock.Mock
}

func (bus *MockEventBus) Subscribe() error {
	ret := bus.Called()
	return ret.Error(0)
}

func (bus *MockEventBus) Publish(event Event) {
	bus.Called(event)
}

func (bus *MockEventBus) PublishAsync(cxt context.Context, event Event) {
	bus.Called(event)
}
