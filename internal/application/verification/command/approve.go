package command

import (
	"context"
	"fmt"

	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/service"
)

const ApproveVerificationCommandType bus.CommandType = "approve.verification.command"

// ApproveVerificationCommand is the command dispatched to approve verification.
type ApproveVerificationCommand struct {
	uuid string
}

// NewApproveVerificationCommand creates a new ApproveVerificationCommand.
func NewApproveVerificationCommand(uuid string) ApproveVerificationCommand {
	return ApproveVerificationCommand{
		uuid: uuid,
	}
}

// Type implements bus.Command interface.
func (c ApproveVerificationCommand) Type() bus.CommandType {
	return ApproveVerificationCommandType
}

// ApproveVerificationCommandHandler is the ApproveVerificationCommand handler.
type ApproveVerificationCommandHandler struct {
	approveVerificationService service.ApproveVerificationService
}

// NewApproveVerificationCommandHandler initializes a new ApproveVerificationCommandHandler.
func NewApproveVerificationCommandHandler(approveVerificationService service.ApproveVerificationService) ApproveVerificationCommandHandler {
	return ApproveVerificationCommandHandler{
		approveVerificationService: approveVerificationService,
	}
}

// Handle implements the bus.CommandHandler interface.
func (h ApproveVerificationCommandHandler) Handle(ctx context.Context, cmd bus.Command) error {
	approveVerificationCommand, ok := cmd.(ApproveVerificationCommand)
	if !ok {
		return fmt.Errorf("command type %s: %w", cmd.Type(), bus.ErrUnexpectedCommand)
	}

	return h.approveVerificationService.Approve(ctx, approveVerificationCommand.uuid)
}
