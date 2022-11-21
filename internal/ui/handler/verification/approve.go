package verification

import (
	"context"
	"net/http"

	"github.com/vitalii-tkachuk/verification-service/internal/application/verification/command"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure"
)

// approveVerificationResponse represents approve verification endpoint response structure.
type approveVerificationResponse struct {
	Uuid string `json:"uuid"`
}

// ApproveVerificationHandler returns an HTTP handler for verification approval.
func ApproveVerificationHandler(application *infrastructure.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		verificationUuid := application.GetUrlParam(r, "verificationUuid")

		approveCommand := command.NewApproveVerificationCommand(verificationUuid)

		if err := application.CommandBus.Dispatch(context.Background(), approveCommand); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}

		response := approveVerificationResponse{Uuid: verificationUuid}

		if err := application.Marshall(w, http.StatusOK, response, nil); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}
	}
}
