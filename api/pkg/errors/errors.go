package errors

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func ValidationErrorsToJSON(err error) ErrorResponse {
	var errs validator.ValidationErrors
	errors.As(err, &errs)
	errorResponse := ErrorResponse{
		Errors: make(map[string]string),
	}

	for _, e := range errs {
		errorResponse.Errors[e.Field()] = e.Tag()
	}

	return errorResponse
}
