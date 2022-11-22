package command

import (
	"context"
	"fmt"

	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/service"
)

const DeclineVerificationCommandType bus.CommandType = "decline.verification.command"

// DeclineVerificationCommand is the command dispatched to decline verification.
type DeclineVerificationCommand struct {
	uuid, declineReason string
}

// NewDeclineVerificationCommand creates a new DeclineVerificationCommand.
func NewDeclineVerificationCommand(UUID, declineReason string) DeclineVerificationCommand {
	return DeclineVerificationCommand{
		uuid:          UUID,
		declineReason: declineReason,
	}
}

// Type implements bus.Command interface.
func (c DeclineVerificationCommand) Type() bus.CommandType {
	return DeclineVerificationCommandType
}

// DeclineVerificationCommandHandler is the DeclineVerificationCommand handler.
type DeclineVerificationCommandHandler struct {
	declineVerificationService service.DeclineVerificationService
}

// NewDeclineVerificationCommandHandler initializes a new DeclineVerificationCommandHandler.
func NewDeclineVerificationCommandHandler(declineVerificationService service.DeclineVerificationService) DeclineVerificationCommandHandler {
	return DeclineVerificationCommandHandler{
		declineVerificationService: declineVerificationService,
	}
}

// Handle implements the bus.CommandHandler interface.
func (h DeclineVerificationCommandHandler) Handle(ctx context.Context, cmd bus.Command) error {
	declineVerificationCommand, ok := cmd.(DeclineVerificationCommand)
	if !ok {
		return fmt.Errorf("command type %s: %w", cmd.Type(), bus.ErrUnexpectedCommand)
	}

	return h.declineVerificationService.Decline(
		ctx,
		declineVerificationCommand.uuid,
		declineVerificationCommand.declineReason,
	)
}
