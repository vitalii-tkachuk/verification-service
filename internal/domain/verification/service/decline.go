package service

import (
	"context"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
)

// DeclineVerificationService is the default Verification decline service.
type DeclineVerificationService struct {
	verificationRepository aggregate.VerificationRepository
}

// NewDeclineVerificationService returns the default DeclineVerificationService interface implementation.
func NewDeclineVerificationService(verificationRepository aggregate.VerificationRepository) DeclineVerificationService {
	return DeclineVerificationService{
		verificationRepository: verificationRepository,
	}
}

// Decline implements the DeclineVerificationService interface
func (s DeclineVerificationService) Decline(ctx context.Context, uuid, declineReason string) error {
	verificationUuid, err := aggregate.NewVerificationUuid(uuid)
	if err != nil {
		return err
	}

	verification, err := s.verificationRepository.GetByUuid(ctx, verificationUuid)
	if err != nil {
		return err
	}

	if err := verification.Decline(declineReason); err != nil {
		return err
	}

	return s.verificationRepository.Update(ctx, verification)
}
