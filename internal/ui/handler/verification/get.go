package verification

import (
	"net/http"
	"time"

	"github.com/vitalii-tkachuk/verification-service/internal/application/verification/query"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure"
)

// getVerificationByUUIDResponse represents get verification by uuid endpoint response structure.
type getVerificationByUUIDResponse struct {
	ID            uint32    `json:"id"`
	UUID          string    `json:"uuid"`
	Kind          string    `json:"kind"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	DeclineReason string    `json:"declineReason,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

// toVerificationByUUIDResponse create getVerificationByUUIDResponse from aggregate.Verification.
func toVerificationByUUIDResponse(verification *aggregate.Verification) *getVerificationByUUIDResponse {
	return &getVerificationByUUIDResponse{
		ID:            verification.ID().Value(),
		UUID:          verification.UUID().Value(),
		Kind:          verification.Kind().Value(),
		Description:   verification.Description().Value(),
		Status:        verification.Status().Value(),
		DeclineReason: verification.DeclineReason().Value(),
		CreatedAt:     verification.CreatedAt(),
	}
}

// GetVerificationHandler returns an HTTP handler for verification fetching.
func GetVerificationHandler(application *infrastructure.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		verificationUUID := application.GetURLParam(r, "verificationUuid")

		getVerificationByUUIDQuery := query.NewGetVerificationByUUIDQuery(verificationUUID)

		verification, err := application.QueryBus.Ask(r.Context(), getVerificationByUUIDQuery)
		if err != nil {
			application.HttpErrorResponse(w, err)

			return
		}

		response := toVerificationByUUIDResponse(verification.(*aggregate.Verification))

		if err := application.Marshall(w, http.StatusOK, response, nil); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}
	}
}
