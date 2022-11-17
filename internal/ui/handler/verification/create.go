package verification

import (
	"context"
	"github.com/google/uuid"
	"github.com/vitalii-tkachuk/verification-service/internal/application/verification/command"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure"
	"net/http"
)

// createVerificationRequest represents create verification endpoint structure.
type createVerificationRequest struct {
	Description string `json:"description" validate:"required,min=10"`
	Kind        string `json:"kind" validate:"required,alpha,oneof=identity document"`
}

// createVerificationResponse represents create verification endpoint response structure.
type createVerificationResponse struct {
	Uuid uuid.UUID `json:"uuid"`
}

// CreateVerificationHandler returns an HTTP handler for verification creation.
func CreateVerificationHandler(application *infrastructure.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request createVerificationRequest

		if err := application.Unmarshall(w, r, &request); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}

		if err := application.ValidateRequest(request); err != nil {
			application.ValidationErrorResponse(w, err)

			return
		}

		verificationUuid := uuid.New()
		createCommand := command.NewCreateVerificationCommand(verificationUuid, request.Description, request.Kind)

		if err := application.CommandBus.Dispatch(context.Background(), createCommand); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}

		response := createVerificationResponse{Uuid: verificationUuid}

		if err := application.Marshall(w, http.StatusCreated, response, nil); err != nil {
			application.HttpErrorResponse(w, err)

			return
		}
	}
}
