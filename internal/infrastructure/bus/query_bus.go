package bus

import (
	"context"
	"fmt"

	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
)

// QueryBus represents in memory query bus implementation of bus.QueryBus interface.
type QueryBus struct {
	handlers map[bus.QueryType]bus.QueryHandler
}

// NewQueryBus creates a new QueryBus.
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[bus.QueryType]bus.QueryHandler),
	}
}

// Ask implements bus.CommandBus.Ask method.
func (b QueryBus) Ask(ctx context.Context, query bus.Query) (interface{}, error) {
	handler, ok := b.handlers[query.Type()]
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", query.Type(), bus.ErrQueryHandlerNotFound)
	}

	return handler.Handle(ctx, query)
}

// Register implements bus.CommandBus.Register method.
func (b QueryBus) Register(queryType bus.QueryType, handler bus.QueryHandler) {
	b.handlers[queryType] = handler
}
