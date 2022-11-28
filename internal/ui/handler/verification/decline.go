package verification

import (
	"net/http"

	"github.com/vitalii-tkachuk/verification-service/internal/application/verification/command"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure"
)

// declineVerificationResponse represents decline verification endpoint response structure.
type declineVerificationRequest struct {
	DeclineReason string `json:"declineReason" validate:"required,min=5"`
}

// declineVerificationResponse represents decline verification endpoint response structure.
type declineVerificationResponse struct {
	UUID string `json:"uuid"`
}

// DeclineVerificationHandler returns an HTTP handler for verification decline.
func DeclineVerificationHandler(application *infrastructure.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request declineVerificationRequest

		if err := application.Unmarshall(w, r, &request); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}

		if err := application.ValidateRequest(request); err != nil {
			application.ValidationErrorResponse(w, err)

			return
		}

		verificationUUID := application.GetURLParam(r, "verificationUuid")
		declineCommand := command.NewDeclineVerificationCommand(verificationUUID, request.DeclineReason)

		if err := application.CommandBus.Dispatch(r.Context(), declineCommand); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}

		response := declineVerificationResponse{UUID: verificationUUID}

		if err := application.Marshall(w, http.StatusOK, response, nil); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}
	}
}
