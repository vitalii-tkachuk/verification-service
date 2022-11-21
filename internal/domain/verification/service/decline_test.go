package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/persistence/postgres"
	"github.com/vitalii-tkachuk/verification-service/test/mocks/persistence"
)

func TestDeclineVerificationServiceInvalidUuidError(t *testing.T) {
	// assign
	verificationUuid := "invalidUuid"
	declineReason := "Bad document quantity"

	// act
	verificationRepositoryMock := new(persistence.VerificationRepository)
	declineVerificationService := NewDeclineVerificationService(verificationRepositoryMock)
	err := declineVerificationService.Decline(context.Background(), verificationUuid, declineReason)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, aggregate.ErrInvalidVerificationUuid)
}

func TestDeclineVerificationServiceNotFoundError(t *testing.T) {
	// assign
	verificationUuid := uuid.New()
	declineReason := "Bad document quantity"

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUuid", mock.Anything, mock.Anything).Return(nil, postgres.ErrVerificationNotFound)

	// act
	declineVerificationService := NewDeclineVerificationService(verificationRepositoryMock)
	err := declineVerificationService.Decline(context.Background(), verificationUuid.String(), declineReason)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, postgres.ErrVerificationNotFound)
}

func TestDeclineVerificationServiceAlreadyProcessedError(t *testing.T) {
	// assign
	declineReason := "Bad document quantity"
	processedVerification, _ := aggregate.NewVerification(
		uuid.New().String(),
		aggregate.Identity,
		"Fancy verification document description",
	)
	_ = processedVerification.Approve()

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUuid", mock.Anything, mock.Anything).Return(processedVerification, nil)

	// act
	declineVerificationService := NewDeclineVerificationService(verificationRepositoryMock)
	err := declineVerificationService.Decline(context.Background(), processedVerification.Uuid().Value(), declineReason)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, aggregate.ErrAlreadyProcessed)
}

func TestDeclineVerificationServiceSuccess(t *testing.T) {
	// assign
	declineReason := "Bad document quantity"
	verification, _ := aggregate.NewVerification(
		uuid.New().String(),
		aggregate.Identity,
		"Fancy verification document description",
	)

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUuid", mock.Anything, mock.Anything).Return(verification, nil)
	verificationRepositoryMock.On("Update", mock.Anything, mock.Anything).Return(nil)

	// act
	declineVerificationService := NewDeclineVerificationService(verificationRepositoryMock)
	err := declineVerificationService.Decline(context.Background(), verification.Uuid().Value(), declineReason)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, aggregate.Declined, verification.Status().Value())
}
