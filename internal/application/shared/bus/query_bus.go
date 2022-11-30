package bus

import (
	"context"
	"errors"
)

var (
	ErrUnexpectedQuery      = errors.New("unexpected query")
	ErrQueryHandlerNotFound = errors.New("query handler not found")
)

// QueryBus defines interface for CQRS query bus implementations.
type QueryBus interface {
	Ask(context.Context, Query) (any, error)
	Register(QueryType, QueryHandler)
}

// QueryType is a unique string needed to identity query in QueryBus.
type QueryType string

// Query defines interface for query in QueryBus with unique identifier.
type Query interface {
	Type() QueryType
}

//go:generate mockery --case=snake --outpkg=mocks --output=test/mocks --name=Query

// QueryHandler defined interface for query handler. Handler match command by command type.
type QueryHandler interface {
	Handle(context.Context, Query) (any, error)
}
