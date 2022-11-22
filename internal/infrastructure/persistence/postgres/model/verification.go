package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
)

const (
	SQLVerificationTable     = "verifications"
	SQLVerificationCreateTag = "create"
	SQLVerificationGetTag    = "get"
)

var ErrFailedRestoringVerificationFromDatabase = errors.New("failed restoring verification from database")

// SQLVerification verification represents aggregate.Verification database structure.
// separate struct is used because aggregate with VO is hard to persist to database
type SQLVerification struct {
	ID            uint32    `db:"id" fieldtag:"get"`
	UUID          string    `db:"uuid" fieldtag:"create,get"`
	Kind          string    `db:"kind" fieldtag:"create,get"`
	Description   string    `db:"description" fieldtag:"create,get"`
	Status        string    `db:"status" fieldtag:"create,get"`
	DeclineReason string    `db:"decline_reason" fieldtag:"create,get"`
	CreatedAt     time.Time `db:"created_at" fieldtag:"create,get"`
}

// ToSQLVerification convert aggregate.Verification to it's sql representation.
func ToSQLVerification(verification *aggregate.Verification) SQLVerification {
	sqlVerification := SQLVerification{
		UUID:        verification.UUID().Value(),
		Kind:        verification.Kind().Value(),
		Description: verification.Description().Value(),
		Status:      verification.Status().Value(),
		CreatedAt:   verification.CreatedAt(),
	}

	if verification.DeclineReason().Value() != "" {
		sqlVerification.DeclineReason = verification.DeclineReason().Value()
	}

	return sqlVerification
}

// ToDomainVerification convert SqlVerification to domain aggregate.
func ToDomainVerification(sqlVerification SQLVerification) (*aggregate.Verification, error) {
	verification, err := aggregate.NewVerification(
		sqlVerification.UUID,
		sqlVerification.Kind,
		sqlVerification.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrFailedRestoringVerificationFromDatabase, err)
	}

	verification.WithID(sqlVerification.ID)

	if err = verification.WithStatus(sqlVerification.Status); err != nil {
		return nil, err
	}

	err = verification.WithDeclineReason(sqlVerification.DeclineReason)

	if sqlVerification.DeclineReason != "" && err != nil {
		return nil, err
	}

	return verification, nil
}
