package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/persistence/postgres/model"
	"time"
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
	verificationSQLStruct := sqlbuilder.NewStruct(new(model.SqlVerification))

	selectBuilder := verificationSQLStruct.InsertIntoForTag(
		model.SqlVerificationTable,
		model.SqlVerificationCreateTag,
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
	verificationSQLStruct := sqlbuilder.NewStruct(new(model.SqlVerification))

	updateBuilder := verificationSQLStruct.UpdateForTag(
		model.SqlVerificationTable,
		model.SqlVerificationCreateTag,
		model.ToSQLVerification(verification),
	)
	updateBuilder.Where(updateBuilder.Equal("uuid", verification.Uuid().Value()))

	query, args := updateBuilder.BuildWithFlavor(sqlbuilder.PostgreSQL)

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(ctxTimeout, query, args...); err != nil {
		return fmt.Errorf("%s: %w", ErrVerificationPersistFailed, err)
	}

	return nil
}

// GetByUuid implements the aggregate.VerificationRepository.GetByUuid() method.
func (r *VerificationRepository) GetByUuid(ctx context.Context, uuid aggregate.VerificationUuid) (*aggregate.Verification, error) {
	verificationSQLStruct := sqlbuilder.NewStruct(new(model.SqlVerification))

	selectBuilder := verificationSQLStruct.SelectFromForTag(model.SqlVerificationTable, model.SqlVerificationGetTag)
	selectBuilder.Where(selectBuilder.Equal("uuid", uuid.Value()))

	query, args := selectBuilder.BuildWithFlavor(sqlbuilder.PostgreSQL)

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var sqlVerification model.SqlVerification

	err := r.db.QueryRowContext(ctxTimeout, query, args...).Scan(verificationSQLStruct.Addr(&sqlVerification)...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: %s", ErrVerificationNotFound, uuid.Value())
		default:
			return nil, err
		}
	}

	return model.ToDomainVerification(sqlVerification)
}
