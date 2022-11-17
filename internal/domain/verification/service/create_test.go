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

func TestCreateVerificationServiceDomainError(t *testing.T) {
	// assign
	verificationUuid := uuid.New()
	description := ""
	kind := aggregate.Identity

	// act
	verificationRepositoryMock := new(persistence.VerificationRepository)
	createVerificationService := NewCreateVerificationService(verificationRepositoryMock)
	err := createVerificationService.Create(context.Background(), verificationUuid, description, kind)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, aggregate.ErrEmptyDescription)
}

func TestCreateVerificationServicePersistenceError(t *testing.T) {
	// assign
	verificationUuid := uuid.New()
	description := "Fancy verification document description"
	kind := aggregate.Document

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("Add", mock.Anything, mock.Anything).Return(postgres.ErrVerificationPersistFailed)

	// act
	createVerificationService := NewCreateVerificationService(verificationRepositoryMock)
	err := createVerificationService.Create(context.Background(), verificationUuid, description, kind)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, postgres.ErrVerificationPersistFailed)
}

func TestCreateVerificationServiceSuccess(t *testing.T) {
	// assign
	verificationUuid := uuid.New()
	description := "Fancy verification document description"
	kind := aggregate.Identity

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("Add", mock.Anything, mock.Anything).Return(nil)

	// act
	createVerificationService := NewCreateVerificationService(verificationRepositoryMock)
	err := createVerificationService.Create(context.Background(), verificationUuid, description, kind)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
}
