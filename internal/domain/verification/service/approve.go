package service

import (
	"context"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
)

// ApproveVerificationService is the default Verification approve service
type ApproveVerificationService struct {
	verificationRepository aggregate.VerificationRepository
}

// NewApproveVerificationService returns the default CreateVerificationService interface implementation
func NewApproveVerificationService(verificationRepository aggregate.VerificationRepository) ApproveVerificationService {
	return ApproveVerificationService{
		verificationRepository: verificationRepository,
	}
}

// Approve implements the ApproveVerificationService interface
func (s ApproveVerificationService) Approve(ctx context.Context, uuid string) error {
	verificationUuid, err := aggregate.NewVerificationUuid(uuid)
	if err != nil {
		return err
	}

	verification, err := s.verificationRepository.GetByUuid(ctx, verificationUuid)
	if err != nil {
		return err
	}

	if err := verification.Approve(); err != nil {
		return err
	}

	return s.verificationRepository.Update(ctx, verification)
}
