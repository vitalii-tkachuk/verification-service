package query

import (
	"context"
	"fmt"

	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
)

const GetVerificationByUuidQueryType bus.QueryType = "get_by_uuid.verification.query"

// GetVerificationByUuidQuery is the query dispatched to get verification by uuid.
type GetVerificationByUuidQuery struct {
	uuid string
}

// NewGetVerificationByUuidQuery creates a new GetVerificationByUuidQuery.
func NewGetVerificationByUuidQuery(uuid string) GetVerificationByUuidQuery {
	return GetVerificationByUuidQuery{
		uuid: uuid,
	}
}

// Type implements bus.Query interface.
func (q GetVerificationByUuidQuery) Type() bus.QueryType {
	return GetVerificationByUuidQueryType
}

// GetVerificationByUuidQueryHandler is the GetVerificationByUuidQuery handler.
type GetVerificationByUuidQueryHandler struct {
	verificationRepository aggregate.VerificationRepository
}

// NewGetVerificationByUuidQueryHandler initializes a new GetVerificationByUuidQueryHandler.
func NewGetVerificationByUuidQueryHandler(verificationRepository aggregate.VerificationRepository) GetVerificationByUuidQueryHandler {
	return GetVerificationByUuidQueryHandler{
		verificationRepository: verificationRepository,
	}
}

// Handle implements the bus.QueryHandler interface.
func (h GetVerificationByUuidQueryHandler) Handle(ctx context.Context, q bus.Query) (interface{}, error) {
	getVerificationCommand, ok := q.(GetVerificationByUuidQuery)
	if !ok {
		return nil, fmt.Errorf("query type %s: %w", q.Type(), bus.ErrUnexpectedQuery)
	}

	verificationUuid, err := aggregate.NewVerificationUuid(getVerificationCommand.uuid)
	if err != nil {
		return nil, err
	}

	return h.verificationRepository.GetByUuid(ctx, verificationUuid)
}
