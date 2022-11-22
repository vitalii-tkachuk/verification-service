package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/service"
)

const CreateVerificationCommandType bus.CommandType = "create.verification.command"

// CreateVerificationCommand is the command dispatched to create a new verification
type CreateVerificationCommand struct {
	uuid        uuid.UUID
	description string
	kind        string
}

// NewCreateVerificationCommand creates a new CreateVerificationCommand
func NewCreateVerificationCommand(UUID uuid.UUID, description, kind string) CreateVerificationCommand {
	return CreateVerificationCommand{
		uuid:        UUID,
		description: description,
		kind:        kind,
	}
}

// Type implements bus.Command interface
func (c CreateVerificationCommand) Type() bus.CommandType {
	return CreateVerificationCommandType
}

// CreateVerificationCommandHandler is the CreateVerificationCommand handler
type CreateVerificationCommandHandler struct {
	createVerificationService service.CreateVerificationService
}

// NewCreateVerificationCommandHandler initializes a new CreateVerificationCommandHandler.
func NewCreateVerificationCommandHandler(createVerificationService service.CreateVerificationService) CreateVerificationCommandHandler {
	return CreateVerificationCommandHandler{
		createVerificationService: createVerificationService,
	}
}

// Handle implements the bus.CommandHandler interface.
func (h CreateVerificationCommandHandler) Handle(ctx context.Context, cmd bus.Command) error {
	createVerificationCommand, ok := cmd.(CreateVerificationCommand)
	if !ok {
		return fmt.Errorf("command type %s: %w", cmd.Type(), bus.ErrUnexpectedCommand)
	}

	return h.createVerificationService.Create(
		ctx,
		createVerificationCommand.uuid,
		createVerificationCommand.description,
		createVerificationCommand.kind,
	)
}
