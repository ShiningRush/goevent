package goevent

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockEventBus used to execute unit test
type MockEventBus struct {
	mock.Mock
}

func (bus *MockEventBus) Subscribe(handler EventHandler) error {
	ret := bus.Called(handler)
	return ret.Error(0)
}

func (bus *MockEventBus) Publish(event Event) {
	bus.Called(event)
}

func (bus *MockEventBus) PublishSync(ctx context.Context, event Event) {
	bus.Called(ctx, event)
}

func (bus *MockEventBus) Close() {
	bus.Called()
}
