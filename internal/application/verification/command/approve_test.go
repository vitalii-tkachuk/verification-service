package command

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/service"
	"github.com/vitalii-tkachuk/verification-service/test/mocks"
	"github.com/vitalii-tkachuk/verification-service/test/mocks/persistence"
	"testing"
)

func TestHandleUnsupportedCommandrror(t *testing.T) {
	// assign
	var unsupportedCommandType bus.CommandType = "unsupported_approve.verification.command"

	unsupportedCommand := new(mocks.Command)
	unsupportedCommand.On("Type").Return(unsupportedCommandType)
	verificationRepositoryMock := new(persistence.VerificationRepository)

	// act
	approveVerificationService := service.NewApproveVerificationService(verificationRepositoryMock)
	approveVerificationCommandHandler := NewApproveVerificationCommandHandler(approveVerificationService)
	err := approveVerificationCommandHandler.Handle(context.Background(), unsupportedCommand)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.ErrorIs(t, err, bus.ErrUnexpectedCommand)
}

func TestHandleApproveVerificationCommandSuccess(t *testing.T) {
	// assign
	verification, _ := aggregate.NewVerification(
		uuid.New().String(),
		aggregate.Identity,
		"Fancy verification document description",
	)

	approveVerificationCommand := NewApproveVerificationCommand(verification.Uuid().Value())

	verificationRepositoryMock := new(persistence.VerificationRepository)
	verificationRepositoryMock.On("GetByUuid", mock.Anything, mock.Anything).Return(verification, nil)
	verificationRepositoryMock.On("Update", mock.Anything, mock.Anything).Return(nil)

	// act
	approveVerificationService := service.NewApproveVerificationService(verificationRepositoryMock)

	approveVerificationCommandHandler := NewApproveVerificationCommandHandler(approveVerificationService)
	err := approveVerificationCommandHandler.Handle(context.Background(), approveVerificationCommand)

	// assert
	verificationRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, aggregate.Approved, verification.Status().Value())
}
