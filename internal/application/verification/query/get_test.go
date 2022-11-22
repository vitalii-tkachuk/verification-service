package query

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/test/mocks"
	"github.com/vitalii-tkachuk/verification-service/test/mocks/persistence"
)

func TestHandleUnsupportedGetByUUIDQueryError(t *testing.T) {
	// assign
	var unsupportedQueryType bus.QueryType = "unsupported_get_by_uuid.verification.query"

	unsupportedQuery := new(mocks.Query)
	unsupportedQuery.On("Type").Return(unsupportedQueryType)

	verificationRepositoryMock := new(persistence.VerificationRepository)

	// act
	getVerificationByUUIDQueryHandler := NewGetVerificationByUUIDQueryHandler(verificationRepositoryMock)
	verification, err := getVerificationByUUIDQueryHandler.Handle(context.Background(), unsupportedQuery)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.Nil(t, verification)
	assert.ErrorIs(t, err, bus.ErrUnexpectedQuery)
}

func TestGetVerificationByUUIDQuerySuccess(t *testing.T) {
	// assign
	expectedVerification, _ := aggregate.NewVerification(
		uuid.New().String(),
		aggregate.Identity,
		"Fancy verification document description",
	)

	getVerificationByUUIDQuery := NewGetVerificationByUUIDQuery(expectedVerification.UUID().Value())

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUUID", mock.Anything, mock.Anything).Return(expectedVerification, nil)

	// act
	getVerificationByUUIDQueryHandler := NewGetVerificationByUUIDQueryHandler(verificationRepositoryMock)
	verification, err := getVerificationByUUIDQueryHandler.Handle(context.Background(), getVerificationByUUIDQuery)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, expectedVerification.UUID(), verification.(*aggregate.Verification).UUID())
}
