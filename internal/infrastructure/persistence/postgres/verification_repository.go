package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/persistence/postgres/model"
)

var (
	ErrVerificationPersistFailed = errors.New("error trying to persist verification to database")
	ErrVerificationNotFound      = errors.New("verification not found")
)

// VerificationRepository is a PostgreSQL aggregate.VerificationRepository implementation.
type VerificationRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewVerificationRepository initializes a PostgreSQL-based implementation of aggregate.VerificationRepository.
func NewVerificationRepository(db *sql.DB, dbTimeout time.Duration) *VerificationRepository {
	return &VerificationRepository{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

// Add implements the aggregate.VerificationRepository.Add() method.
func (r *VerificationRepository) Add(ctx context.Context, verification *aggregate.Verification) error {
	verificationSQLStruct := sqlbuilder.NewStruct(new(model.SQLVerification))

	selectBuilder := verificationSQLStruct.InsertIntoForTag(
		model.SQLVerificationTable,
		model.SQLVerificationCreateTag,
		model.ToSQLVerification(verification),
	)
	query, args := selectBuilder.BuildWithFlavor(sqlbuilder.PostgreSQL)

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(ctxTimeout, query, args...); err != nil {
		return fmt.Errorf("%s: %w", ErrVerificationPersistFailed, err)
	}

	return nil
}

// Update implements the aggregate.VerificationRepository.Update() method.
func (r *VerificationRepository) Update(ctx context.Context, verification *aggregate.Verification) error {
	verificationSQLStruct := sqlbuilder.NewStruct(new(model.SQLVerification))

	updateBuilder := verificationSQLStruct.UpdateForTag(
		model.SQLVerificationTable,
		model.SQLVerificationCreateTag,
		model.ToSQLVerification(verification),
	)
	updateBuilder.Where(updateBuilder.Equal("uuid", verification.UUID().Value()))

	query, args := updateBuilder.BuildWithFlavor(sqlbuilder.PostgreSQL)

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(ctxTimeout, query, args...); err != nil {
		return fmt.Errorf("%s: %w", ErrVerificationPersistFailed, err)
	}

	return nil
}

// GetByUUID implements the aggregate.VerificationRepository.GetByUUID() method.
func (r *VerificationRepository) GetByUUID(ctx context.Context, uuid aggregate.VerificationUUID) (*aggregate.Verification, error) {
	verificationSQLStruct := sqlbuilder.NewStruct(new(model.SQLVerification))

	selectBuilder := verificationSQLStruct.SelectFromForTag(model.SQLVerificationTable, model.SQLVerificationGetTag)
	selectBuilder.Where(selectBuilder.Equal("uuid", uuid.Value()))

	query, args := selectBuilder.BuildWithFlavor(sqlbuilder.PostgreSQL)

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var SQLVerification model.SQLVerification

	err := r.db.QueryRowContext(ctxTimeout, query, args...).Scan(verificationSQLStruct.Addr(&SQLVerification)...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: %s", ErrVerificationNotFound, uuid.Value())
		default:
			return nil, err
		}
	}

	return model.ToDomainVerification(SQLVerification)
}
