package model

import (
	"errors"
	"fmt"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"time"
)

const (
	SqlVerificationTable     = "verifications"
	SqlVerificationCreateTag = "create"
	SqlVerificationGetTag    = "get"
)

var ErrFailedRestoringVerificationFromDatabase = errors.New("failed restoring verification from database")

// SqlVerification verification represents aggregate.Verification database structure.
// separate struct is used because aggregate with VO is hard to persist to database
type SqlVerification struct {
	Id            uint32    `db:"id" fieldtag:"get"`
	Uuid          string    `db:"uuid" fieldtag:"create,get"`
	Kind          string    `db:"kind" fieldtag:"create,get"`
	Description   string    `db:"description" fieldtag:"create,get"`
	Status        string    `db:"status" fieldtag:"create,get"`
	DeclineReason string    `db:"decline_reason" fieldtag:"create,get"`
	CreatedAt     time.Time `db:"created_at" fieldtag:"create,get"`
}

// ToSQLVerification convert aggregate.Verification to it's sql representation.
func ToSQLVerification(verification *aggregate.Verification) SqlVerification {
	sqlVerification := SqlVerification{
		Uuid:        verification.Uuid().Value(),
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
func ToDomainVerification(sqlVerification SqlVerification) (*aggregate.Verification, error) {
	verification, err := aggregate.NewVerification(
		sqlVerification.Uuid,
		sqlVerification.Kind,
		sqlVerification.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrFailedRestoringVerificationFromDatabase, err)
	}

	verification.WithId(sqlVerification.Id)

	if err = verification.WithStatus(sqlVerification.Status); err != nil {
		return nil, err
	}

	err = verification.WithDeclineReason(sqlVerification.DeclineReason)

	if sqlVerification.DeclineReason != "" && err != nil {
		return nil, err
	}

	return verification, nil
}
