package verification

import (
	"net/http"

	"github.com/vitalii-tkachuk/verification-service/internal/application/verification/command"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure"
)

// approveVerificationResponse represents approve verification endpoint response structure.
type approveVerificationResponse struct {
	UUID string `json:"uuid"`
}

// ApproveVerificationHandler returns an HTTP handler for verification approval.
func ApproveVerificationHandler(application *infrastructure.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		verificationUUID := application.GetURLParam(r, "verificationUuid")

		approveCommand := command.NewApproveVerificationCommand(verificationUUID)

		if err := application.CommandBus.Dispatch(r.Context(), approveCommand); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}

		response := approveVerificationResponse{UUID: verificationUUID}

		if err := application.Marshall(w, http.StatusOK, response, nil); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}
	}
}
