package aggregate

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerification(t *testing.T) {
	t.Parallel()

	t.Run("test create verification id success", testCreateVerificationIdSuccess)
	t.Run("test create verification uuid success", testCreateVerificationUuidSuccess)
	t.Run("test create invalid verification uuid error", testCreateInvalidVerificationUuidError)
	t.Run("test create verification kind success", testCreateVerificationKindSuccess)
	t.Run("test create verification kind error", testCreateInvalidVerificationKindError)
	t.Run("test create verification description success", testCreateVerificationDescriptionSuccess)
	t.Run("test create empty verification description error", testCreateEmptyVerificationDescriptionError)
	t.Run("test create verification status success", testCreateVerificationStatusSuccess)
	t.Run("test create verification status error", testCreateInvalidVerificationStatusError)
	t.Run("test create verification decline reason success", testCreateVerificationDeclineReasonSuccess)
	t.Run("test create empty verification decline reason error", testCreateEmptyVerificationDeclineReasonError)
	t.Run("test create verification success", testCreateVerificationSuccess)
	t.Run("test decline verification success", testDeclineVerificationSuccess)
	t.Run("test decline already processed verification error", testDeclineAlreadyProcessedVerificationError)
	t.Run("test decline verification with empty error", testDeclineVerificationWithEmptyReasonError)
	t.Run("test approve verification success", testApproveVerificationSuccess)
	t.Run("test approve already processed verification error", testApproveAlreadyProcessedVerificationError)
}

func testCreateVerificationIdSuccess(t *testing.T) {
	// assign
	var value uint32 = 100

	// act
	verificationId := NewVerificationId(value)

	// assert
	require.Equal(t, value, verificationId.Value())
}

func testCreateVerificationUuidSuccess(t *testing.T) {
	// assign
	expectedUuid := uuid.New()

	// act
	verificationUuid, err := NewVerificationUuid(expectedUuid.String())

	// assert
	require.NoError(t, err)
	require.Equal(t, expectedUuid.String(), verificationUuid.Value())
}

func testCreateInvalidVerificationUuidError(t *testing.T) {
	// assign
	invalidUuid := "invalidUuid"

	// act
	verificationUuid, err := NewVerificationUuid(invalidUuid)

	// assert
	require.ErrorIs(t, err, ErrInvalidVerificationUuid)
	require.Equal(t, VerificationUuid{}, verificationUuid)
}

func testCreateVerificationKindSuccess(t *testing.T) {
	// assign
	kind := Identity

	// act
	verificationKind, err := NewVerificationKind(kind)

	// assert
	require.NoError(t, err)
	require.Equal(t, kind, verificationKind.Value())
}

func testCreateInvalidVerificationKindError(t *testing.T) {
	// assign
	invalidKind := "invalidVerificationStatus"

	// act
	verificationKind, err := NewVerificationKind(invalidKind)

	// assert
	require.ErrorIs(t, err, ErrInvalidVerificationKind)
	require.Equal(t, VerificationKind{}, verificationKind)
}

func testCreateVerificationDescriptionSuccess(t *testing.T) {
	// assign
	description := "Fancy verification document description"

	// act
	verificationDescription, err := NewVerificationDescription(description)

	// assert
	require.NoError(t, err)
	require.Equal(t, description, verificationDescription.Value())
}

func testCreateEmptyVerificationDescriptionError(t *testing.T) {
	// assign
	description := ""

	// act
	verificationDescription, err := NewVerificationDescription(description)

	// assert
	require.ErrorIs(t, err, ErrEmptyDescription)
	require.Equal(t, VerificationDescription{}, verificationDescription)
}

func testCreateVerificationStatusSuccess(t *testing.T) {
	// assign
	status := Draft

	// act
	verificationStatus, err := NewVerificationStatus(status)

	// assert
	require.NoError(t, err)
	require.Equal(t, status, verificationStatus.Value())
}

func testCreateInvalidVerificationStatusError(t *testing.T) {
	// assign
	status := "invalidStatus"

	// act
	verificationStatus, err := NewVerificationStatus(status)

	// assert
	require.ErrorIs(t, err, ErrInvalidVerificationStatus)
	require.Equal(t, VerificationStatus{}, verificationStatus)
}

func testCreateVerificationDeclineReasonSuccess(t *testing.T) {
	// assign
	declineReason := "Bad photo quality"

	// act
	verificationDeclineReason, err := NewVerificationDeclineReason(declineReason)

	// assert
	require.NoError(t, err)
	require.Equal(t, declineReason, verificationDeclineReason.Value())
}

func testCreateEmptyVerificationDeclineReasonError(t *testing.T) {
	// assign
	declineReason := ""

	// act
	verificationDeclineReason, err := NewVerificationDeclineReason(declineReason)

	// assert
	require.ErrorIs(t, err, ErrEmptyDeclineReason)
	require.Equal(t, VerificationDeclineReason{}, verificationDeclineReason)
}

func testCreateVerificationSuccess(t *testing.T) {
	// assign
	expectedUuid := uuid.New()
	kind := Identity
	description := "Fancy verification document description"

	// act
	verification, err := NewVerification(expectedUuid.String(), kind, description)

	// assert
	require.NoError(t, err)
	require.Equal(t, expectedUuid.String(), verification.Uuid().Value())
	require.Equal(t, kind, verification.Kind().Value())
	require.Equal(t, description, verification.Description().Value())
}

func testDeclineVerificationSuccess(t *testing.T) {
	// assign
	expectedUuid := uuid.New()
	kind := Identity
	description := "Fancy verification document description"
	declineReason := "Bad photo quality"

	// act
	verification, _ := NewVerification(expectedUuid.String(), kind, description)
	err := verification.Decline(declineReason)

	// assert
	require.NoError(t, err)
	require.Equal(t, declineReason, verification.DeclineReason().Value())
	require.Equal(t, Declined, verification.Status().Value())
}

func testDeclineAlreadyProcessedVerificationError(t *testing.T) {
	// assign
	expectedUuid := uuid.New()
	kind := Identity
	description := "Fancy verification document description"
	declineReason := "Bad photo quality"

	// act
	verification, _ := NewVerification(expectedUuid.String(), kind, description)
	_ = verification.Approve()
	err := verification.Decline(declineReason)

	// assert
	require.ErrorIs(t, err, ErrAlreadyProcessed)
	require.Equal(t, Approved, verification.Status().Value())
}

func testDeclineVerificationWithEmptyReasonError(t *testing.T) {
	// assign
	expectedUuid := uuid.New()
	kind := Identity
	description := "Fancy verification document description"
	declineReason := ""

	// act
	verification, _ := NewVerification(expectedUuid.String(), kind, description)
	err := verification.Decline(declineReason)

	// assert
	require.ErrorIs(t, err, ErrEmptyDeclineReason)
	require.Equal(t, Draft, verification.Status().Value())
}

func testApproveVerificationSuccess(t *testing.T) {
	// assign
	expectedUuid := uuid.New()
	kind := Identity
	description := "Fancy verification document description"

	// act
	verification, _ := NewVerification(expectedUuid.String(), kind, description)
	err := verification.Approve()

	// assert
	require.NoError(t, err)
	require.Equal(t, Approved, verification.Status().Value())
}

func testApproveAlreadyProcessedVerificationError(t *testing.T) {
	// assign
	expectedUuid := uuid.New()
	kind := Identity
	description := "Fancy verification document description"
	declineReason := "Bad photo quality"

	// act
	verification, _ := NewVerification(expectedUuid.String(), kind, description)
	_ = verification.Decline(declineReason)
	err := verification.Approve()

	// assert
	require.ErrorIs(t, err, ErrAlreadyProcessed)
	require.Equal(t, Declined, verification.Status().Value())
}
