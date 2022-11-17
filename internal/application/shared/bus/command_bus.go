package bus

import (
	"context"
	"errors"
)

var (
	ErrUnexpectedCommand      = errors.New("unexpected command")
	ErrCommandHandlerNotFound = errors.New("command handler not found")
)

// CommandBus defines interface for CQRS command bus implementations e.g. RabbitMQCommandBus, ApacheKafkaCommandBus.
type CommandBus interface {
	Dispatch(context.Context, Command) error
	Register(CommandType, CommandHandler)
}

// CommandType is a unique string needed to identity command in CommandBus.
type CommandType string

// Command defines interface for command in CommandBus with unique identifier.
type Command interface {
	Type() CommandType
}

//go:generate mockery --case=snake --outpkg=mocks --output=test/mocks --name=Command

// CommandHandler defined interface for command handler. Handler match command by command type.
type CommandHandler interface {
	Handle(context.Context, Command) error
}
