package command

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/service"
	"github.com/vitalii-tkachuk/verification-service/test/mocks"
	"github.com/vitalii-tkachuk/verification-service/test/mocks/persistence"
)

func TestHandleUnsupportedCreateVerificationCommandError(t *testing.T) {
	// assign
	var unsupportedCommandType bus.CommandType = "unsupported_crete.verification.command"

	unsupportedCommand := new(mocks.Command)
	unsupportedCommand.On("Type").Return(unsupportedCommandType)

	verificationRepositoryMock := new(persistence.VerificationRepository)

	// act
	createVerificationService := service.NewCreateVerificationService(verificationRepositoryMock)

	createVerificationCommandHandler := NewCreateVerificationCommandHandler(createVerificationService)
	err := createVerificationCommandHandler.Handle(context.Background(), unsupportedCommand)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, bus.ErrUnexpectedCommand)
}

func TestHandleCreateVerificationCommandSuccess(t *testing.T) {
	// assign
	verificationUUID := uuid.New()
	kind := aggregate.Identity
	description := "Fancy verification document description"

	createVerificationCommand := NewCreateVerificationCommand(verificationUUID, description, kind)

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("Add", mock.Anything, mock.Anything).Return(nil)

	// act
	createVerificationService := service.NewCreateVerificationService(verificationRepositoryMock)

	createVerificationCommandHandler := NewCreateVerificationCommandHandler(createVerificationService)
	err := createVerificationCommandHandler.Handle(context.Background(), createVerificationCommand)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
}
