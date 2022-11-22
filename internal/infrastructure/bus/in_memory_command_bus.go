package bus

import (
	"context"
	"fmt"

	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
)

// InMemoryCommandBus represents in memory command bus implementation of bus.CommandBus interface.
type InMemoryCommandBus struct {
	handlers map[bus.CommandType]bus.CommandHandler
}

// NewInMemoryCommandBus creates a new InMemoryCommandBus.
func NewInMemoryCommandBus() InMemoryCommandBus {
	return InMemoryCommandBus{
		handlers: make(map[bus.CommandType]bus.CommandHandler),
	}
}

// Dispatch implements bus.CommandBus.Dispatch method.
func (b InMemoryCommandBus) Dispatch(ctx context.Context, command bus.Command) error {
	handler, ok := b.handlers[command.Type()]
	if !ok {
		return fmt.Errorf("%s: %w", command.Type(), bus.ErrCommandHandlerNotFound)
	}

	return handler.Handle(ctx, command)
}

// Register implements bus.CommandBus.Register method.
func (b InMemoryCommandBus) Register(commandType bus.CommandType, handler bus.CommandHandler) {
	b.handlers[commandType] = handler
}
