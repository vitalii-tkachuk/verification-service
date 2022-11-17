package aggregate

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/utils"
	"log"
	"time"
)

// VerificationId represents the verification identifier.
type VerificationId struct {
	value uint32
}

// NewVerificationId instantiate the VO for VerificationId.
func NewVerificationId(value uint32) VerificationId {
	return VerificationId{value: value}
}

// Value return the VerificationId value.
func (id VerificationId) Value() uint32 {
	return id.value
}

var ErrInvalidVerificationUuid = errors.New("invalid verification uuid")

// VerificationUuid represents the verification unique identifier.
type VerificationUuid struct {
	value string
}

// NewVerificationUuid instantiate the VO for VerificationUuid.
func NewVerificationUuid(value string) (VerificationUuid, error) {
	if _, err := uuid.Parse(value); err != nil {
		return VerificationUuid{}, fmt.Errorf("%w: %s", ErrInvalidVerificationUuid, value)
	}

	return VerificationUuid{value: value}, nil
}

// Value return the VerificationUuid value.
func (uuid VerificationUuid) Value() string {
	return uuid.value
}

var ErrInvalidVerificationKind = errors.New("invalid verification kind")

const (
	Identity string = "identity"
	Document string = "document"
)

// VerificationKind represents the verification kind.
type VerificationKind struct {
	value string
}

// NewVerificationKind instantiate the VO for VerificationKind
func NewVerificationKind(value string) (VerificationKind, error) {
	if value != Identity && value != Document {
		return VerificationKind{}, ErrInvalidVerificationKind
	}

	return VerificationKind{value: value}, nil
}

// Value return the VerificationKind value.
func (k VerificationKind) Value() string {
	return k.value
}

var ErrEmptyDescription = errors.New("verification description must not be empty")

// VerificationDescription represents the verification description.
type VerificationDescription struct {
	value string
}

// NewVerificationDescription instantiate the VO for VerificationUuid
func NewVerificationDescription(value string) (VerificationDescription, error) {
	if value == "" {
		return VerificationDescription{}, ErrEmptyDescription
	}

	return VerificationDescription{value: value}, nil
}

// Value return the VerificationDescription value.
func (d VerificationDescription) Value() string {
	return d.value
}

var ErrInvalidVerificationStatus = errors.New("invalid verification status")

const (
	Draft    string = "draft"
	Approved string = "approved"
	Declined string = "declined"
)

// VerificationStatus represents the verification status.
type VerificationStatus struct {
	value string
}

// NewVerificationStatus instantiate the VO for VerificationStatus
func NewVerificationStatus(value string) (VerificationStatus, error) {
	if !utils.Contains(value, []string{Draft, Approved, Declined}) {
		return VerificationStatus{}, ErrInvalidVerificationStatus
	}

	return VerificationStatus{value: value}, nil
}

// Value return the NewVerificationStatus value.
func (s VerificationStatus) Value() string {
	return s.value
}

var ErrEmptyDeclineReason = errors.New("verification decline reason must not be empty")

// VerificationDeclineReason represents the verification decline reason.
type VerificationDeclineReason struct {
	value string
}

// NewVerificationDeclineReason instantiate the VO for VerificationDeclineReason
func NewVerificationDeclineReason(value string) (VerificationDeclineReason, error) {
	if value == "" {
		return VerificationDeclineReason{}, ErrEmptyDeclineReason
	}

	return VerificationDeclineReason{value: value}, nil
}

// Value return the NewVerificationStatus value.
func (r VerificationDeclineReason) Value() string {
	return r.value
}

// Verification is the data structure that represents a verification.
type Verification struct {
	id            VerificationId
	uuid          VerificationUuid
	kind          VerificationKind
	description   VerificationDescription
	status        VerificationStatus
	declineReason VerificationDeclineReason
	createdAt     time.Time
}

var ErrAlreadyProcessed = errors.New("verification is already processed")

// VerificationRepository defines the expected behaviour for a verification storage.
type VerificationRepository interface {
	Add(ctx context.Context, verification *Verification) error
	Update(ctx context.Context, verification *Verification) error
	GetByUuid(ctx context.Context, uuid VerificationUuid) (*Verification, error)
}

//go:generate mockery --case=snake --outpkg=persistence --output=test/mocks/persistence --name=VerificationRepository

// NewVerification creates a new verification.
func NewVerification(uuid, kind, description string) (*Verification, error) {
	verificationUuid, err := NewVerificationUuid(uuid)
	if err != nil {
		return nil, err
	}

	verificationKind, err := NewVerificationKind(kind)
	if err != nil {
		return nil, err
	}

	verificationDescription, err := NewVerificationDescription(description)
	if err != nil {
		return nil, err
	}

	verificationStatus, err := NewVerificationStatus(Draft)
	if err != nil {
		return nil, err
	}

	verification := &Verification{
		uuid:        verificationUuid,
		kind:        verificationKind,
		description: verificationDescription,
		status:      verificationStatus,
		createdAt:   time.Now(),
	}

	return verification, nil
}

// WithId add id to verification. Used for restoring object from DB.
func (v *Verification) WithId(id uint32) {
	v.id = NewVerificationId(id)
}

// WithDeclineReason add decline reason to verification. Used for restoring object from DB.
func (v *Verification) WithDeclineReason(reason string) error {
	declineReason, err := NewVerificationDeclineReason(reason)
	if err != nil {
		return err
	}

	v.declineReason = declineReason

	return nil
}

// WithStatus add status to verification. Used for restoring object from DB.
func (v *Verification) WithStatus(status string) error {
	verificationStatus, err := NewVerificationStatus(status)
	if err != nil {
		return err
	}

	v.status = verificationStatus

	return nil
}

// Id returns the Verification identifier.
func (v Verification) Id() VerificationId {
	return v.id
}

// Uuid returns the Verification unique identifier.
func (v Verification) Uuid() VerificationUuid {
	return v.uuid
}

// Kind returns the Verification kind.
func (v Verification) Kind() VerificationKind {
	return v.kind
}

// Description returns the Verification description.
func (v Verification) Description() VerificationDescription {
	return v.description
}

// Status returns the Verification status.
func (v Verification) Status() VerificationStatus {
	return v.status
}

// DeclineReason returns the Verification decline reason.
func (v Verification) DeclineReason() VerificationDeclineReason {
	return v.declineReason
}

// CreatedAt returns the Verification create date.
func (v Verification) CreatedAt() time.Time {
	return v.createdAt
}

// Decline declines Verification with specific reason.
func (v *Verification) Decline(declineReason string) error {
	if v.status.value != Draft {
		return ErrAlreadyProcessed
	}

	verificationDeclineReason, err := NewVerificationDeclineReason(declineReason)
	if err != nil {
		return err
	}

	v.declineReason = verificationDeclineReason

	verificationStatus, err := NewVerificationStatus(Declined)
	if err != nil {
		return err
	}

	v.status = verificationStatus

	return nil
}

// Approve changes Verification status to approved.
func (v *Verification) Approve() error {
	log.Println(v.status.value)
	if v.status.value != Draft {
		return ErrAlreadyProcessed
	}

	verificationStatus, err := NewVerificationStatus(Approved)
	if err != nil {
		return err
	}

	v.status = verificationStatus

	return nil
}
