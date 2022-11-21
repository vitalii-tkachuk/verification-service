package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
)

// CreateVerificationService is the default Verification create service
type CreateVerificationService struct {
	verificationRepository aggregate.VerificationRepository
}

// NewCreateVerificationService returns the default CreateVerificationService interface implementation
func NewCreateVerificationService(verificationRepository aggregate.VerificationRepository) CreateVerificationService {
	return CreateVerificationService{
		verificationRepository: verificationRepository,
	}
}

// Create implements the CreateVerificationService interface
func (s CreateVerificationService) Create(ctx context.Context, uuid uuid.UUID, description, kind string) error {
	verification, err := aggregate.NewVerification(uuid.String(), kind, description)
	if err != nil {
		return err
	}

	return s.verificationRepository.Add(ctx, verification)
}
