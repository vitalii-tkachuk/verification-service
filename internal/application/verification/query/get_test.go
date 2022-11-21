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

func TestHandleUnsupportedGetByUuidQueryError(t *testing.T) {
	// assign
	var unsupportedQueryType bus.QueryType = "unsupported_get_by_uuid.verification.query"

	unsupportedQuery := new(mocks.Query)
	unsupportedQuery.On("Type").Return(unsupportedQueryType)

	verificationRepositoryMock := new(persistence.VerificationRepository)

	// act
	getVerificationByUuidQueryHandler := NewGetVerificationByUuidQueryHandler(verificationRepositoryMock)
	verification, err := getVerificationByUuidQueryHandler.Handle(context.Background(), unsupportedQuery)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.Nil(t, verification)
	assert.ErrorIs(t, err, bus.ErrUnexpectedQuery)
}

func TestGetVerificationByUuidQuerySuccess(t *testing.T) {
	// assign
	expectedVerification, _ := aggregate.NewVerification(
		uuid.New().String(),
		aggregate.Identity,
		"Fancy verification document description",
	)

	getVerificationByUuidQuery := NewGetVerificationByUuidQuery(expectedVerification.Uuid().Value())

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUuid", mock.Anything, mock.Anything).Return(expectedVerification, nil)

	// act
	getVerificationByUuidQueryHandler := NewGetVerificationByUuidQueryHandler(verificationRepositoryMock)
	verification, err := getVerificationByUuidQueryHandler.Handle(context.Background(), getVerificationByUuidQuery)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, expectedVerification.Uuid(), verification.(*aggregate.Verification).Uuid())
}
