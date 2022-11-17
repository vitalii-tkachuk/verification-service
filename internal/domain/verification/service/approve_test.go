package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/persistence/postgres"
	"github.com/vitalii-tkachuk/verification-service/test/mocks/persistence"
	"testing"
)

func TestApproveVerificationServiceInvalidUuidError(t *testing.T) {
	// assign
	verificationUuid := "invalidUuid"

	// act
	verificationRepositoryMock := new(persistence.VerificationRepository)
	approveVerificationService := NewApproveVerificationService(verificationRepositoryMock)
	err := approveVerificationService.Approve(context.Background(), verificationUuid)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, aggregate.ErrInvalidVerificationUuid)
}

func TestApproveVerificationServiceNotFoundError(t *testing.T) {
	// assign
	verificationUuid := uuid.New()

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUuid", mock.Anything, mock.Anything).Return(nil, postgres.ErrVerificationNotFound)

	// act
	approveVerificationService := NewApproveVerificationService(verificationRepositoryMock)
	err := approveVerificationService.Approve(context.Background(), verificationUuid.String())

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, postgres.ErrVerificationNotFound)
}

func TestApproveVerificationServiceAlreadyProcessedError(t *testing.T) {
	// assign
	processedVerification, _ := aggregate.NewVerification(
		uuid.New().String(),
		aggregate.Identity,
		"Fancy verification document description",
	)
	_ = processedVerification.Approve()

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUuid", mock.Anything, mock.Anything).Return(processedVerification, nil)

	// act
	approveVerificationService := NewApproveVerificationService(verificationRepositoryMock)
	err := approveVerificationService.Approve(context.Background(), processedVerification.Uuid().Value())

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, aggregate.ErrAlreadyProcessed)
}

func TestApproveVerificationServiceSuccess(t *testing.T) {
	// assign
	verification, _ := aggregate.NewVerification(
		uuid.New().String(),
		aggregate.Identity,
		"Fancy verification document description",
	)

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUuid", mock.Anything, mock.Anything).Return(verification, nil)
	verificationRepositoryMock.On("Update", mock.Anything, mock.Anything).Return(nil)

	// act
	approveVerificationService := NewApproveVerificationService(verificationRepositoryMock)
	err := approveVerificationService.Approve(context.Background(), verification.Uuid().Value())

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, aggregate.Approved, verification.Status().Value())
}
