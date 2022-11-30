package query

import (
	"context"
	"fmt"

	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
)

const GetVerificationByUUIDQueryType bus.QueryType = "get_by_uuid.verification.query"

// GetVerificationByUUIDQuery is the query dispatched to get verification by uuid.
type GetVerificationByUUIDQuery struct {
	uuid string
}

// NewGetVerificationByUUIDQuery creates a new GetVerificationByUUIDQuery.
func NewGetVerificationByUUIDQuery(UUID string) GetVerificationByUUIDQuery {
	return GetVerificationByUUIDQuery{
		uuid: UUID,
	}
}

// Type implements bus.Query interface.
func (q GetVerificationByUUIDQuery) Type() bus.QueryType {
	return GetVerificationByUUIDQueryType
}

// GetVerificationByUUIDQueryHandler is the GetVerificationByUUIDQuery handler.
type GetVerificationByUUIDQueryHandler struct {
	verificationRepository aggregate.VerificationRepository
}

// NewGetVerificationByUUIDQueryHandler initializes a new GetVerificationByUUIDQueryHandler.
func NewGetVerificationByUUIDQueryHandler(verificationRepository aggregate.VerificationRepository) GetVerificationByUUIDQueryHandler {
	return GetVerificationByUUIDQueryHandler{
		verificationRepository: verificationRepository,
	}
}

// Handle implements the bus.QueryHandler interface.
func (h GetVerificationByUUIDQueryHandler) Handle(ctx context.Context, q bus.Query) (any, error) {
	getVerificationCommand, ok := q.(GetVerificationByUUIDQuery)
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", q.Type(), bus.ErrUnexpectedQuery)
	}

	verificationUUID, err := aggregate.NewVerificationUUID(getVerificationCommand.uuid)
	if err != nil {
		return nil, err
	}

	return h.verificationRepository.GetByUUID(ctx, verificationUUID)
}
