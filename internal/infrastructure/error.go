package infrastructure

// HttpError represents http errors structure.
type HttpError struct {
	Message string `json:"message"`
}

// HttpErrorResponse represents http errors response structure.
type HttpErrorResponse struct {
	Errors []HttpError `json:"errors"`
}

// NewHttpErrorResponse instantiate the HttpErrorResponse.
func NewHttpErrorResponse(errorMessage string) HttpErrorResponse {
	return HttpErrorResponse{Errors: []HttpError{{Message: errorMessage}}}
}

// ValidationError represents incoming request validator error structure.
type ValidationError struct {
	Message      string `json:"message"`
	PropertyPath string `json:"propertyPath"`
}

// ValidationErrorResponse represents validator errors response structure.
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

// NewValidationErrorResponse instantiate the ValidationErrorResponse.
func NewValidationErrorResponse(errors []ValidationError) ValidationErrorResponse {
	return ValidationErrorResponse{Errors: errors}
}
