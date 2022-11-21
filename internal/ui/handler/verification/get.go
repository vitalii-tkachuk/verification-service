package verification

import (
	"context"
	"net/http"
	"time"

	"github.com/vitalii-tkachuk/verification-service/internal/application/verification/query"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/aggregate"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure"
)

// getVerificationByUuidResponse represents get verification by uuid endpoint response structure.
type getVerificationByUuidResponse struct {
	Id            uint32    `json:"id"`
	Uuid          string    `json:"uuid"`
	Kind          string    `json:"kind"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	DeclineReason string    `json:"declineReason,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

// toVerificationByUuidResponse create getVerificationByUuidResponse from aggregate.Verification.
func toVerificationByUuidResponse(verification *aggregate.Verification) *getVerificationByUuidResponse {
	return &getVerificationByUuidResponse{
		Id:            verification.Id().Value(),
		Uuid:          verification.Uuid().Value(),
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
		verificationUuid := application.GetUrlParam(r, "verificationUuid")

		getVerificationByUuidQuery := query.NewGetVerificationByUuidQuery(verificationUuid)

		verification, err := application.QueryBus.Ask(context.Background(), getVerificationByUuidQuery)
		if err != nil {
			application.HttpErrorResponse(w, err)

			return
		}

		response := toVerificationByUuidResponse(verification.(*aggregate.Verification))

		if err := application.Marshall(w, http.StatusOK, response, nil); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}
	}
}
