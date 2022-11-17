package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/vitalii-tkachuk/verification-service/internal/application/shared/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/utils"
	"io"
	"log"
	"net/http"
	"reflect"
)

const requestMaxBodySizeInBytes = 1048576

var (
	ErrMultipleRequestJsonObjects = errors.New("request body must contain a single JSON object")
	ErrInternal                   = errors.New("internal error")
	ErrMarshalFailed              = errors.New("marshal failed")
	ErrEmptyRequestBody           = errors.New("request body must not be empty")
)

// Application represents container for top level services used in handlers.
type Application struct {
	CommandBus bus.CommandBus
	QueryBus   bus.QueryBus
	Validator  *validator.Validate
}

// NewApplication creates a new Application.
func NewApplication(commandBus bus.CommandBus, queryBus bus.QueryBus, validator *validator.Validate) *Application {
	return &Application{
		CommandBus: commandBus,
		QueryBus:   queryBus,
		Validator:  validator,
	}
}

// GetUrlParam returns chi.Router url param
func (a *Application) GetUrlParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

// Marshall serializes response data to json with specific status code.
func (a *Application) Marshall(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	content, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshalling response data of type %s failed: %v", reflect.TypeOf(data), err)
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		return fmt.Errorf("writing response data of type %s failed: %v", reflect.TypeOf(data), err)
	}

	return nil
}

// Unmarshall deserialize request data to destination struct.
func (a *Application) Unmarshall(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, requestMaxBodySizeInBytes)
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(destination); err != nil {
		switch {
		case errors.Is(err, io.EOF):
			return ErrEmptyRequestBody
		default:
			return fmt.Errorf("unmarshalling data of type %s failed: %v", reflect.TypeOf(destination), err)
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return ErrMultipleRequestJsonObjects
	}

	return nil
}

// HttpErrorResponse write error to response with http.StatusBadRequest status code.
func (a *Application) HttpErrorResponse(w http.ResponseWriter, err error) {
	marshalErr := a.Marshall(w, http.StatusBadRequest, NewHttpErrorResponse(err.Error()), nil)

	if marshalErr == nil {
		return
	}

	log.Printf("%s: %s", ErrMarshalFailed, marshalErr)
	_ = a.Marshall(w, http.StatusInternalServerError, NewHttpErrorResponse(ErrInternal.Error()), nil)
}

// ValidateRequest validate request struct with usage of Validator.
func (a *Application) ValidateRequest(s interface{}) error {
	return a.Validator.Struct(s)
}

// ValidationErrorResponse write validation errors to response.
func (a *Application) ValidationErrorResponse(w http.ResponseWriter, err error) {
	var validationErrors []ValidationError

	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{err.Error(), utils.LcFirst(err.Field())})
	}

	_ = a.Marshall(w, http.StatusBadRequest, NewValidationErrorResponse(validationErrors), nil)
}
