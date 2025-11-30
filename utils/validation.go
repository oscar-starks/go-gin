package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResponse represents the validation error response
type ValidationResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

// FormatValidationErrors converts gin validation errors to a readable format
func FormatValidationErrors(err error) ValidationResponse {
	var validationErrors []ValidationError

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range errs {
			field := strings.ToLower(fieldError.Field())
			message := getValidationMessage(fieldError)

			validationErrors = append(validationErrors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}

	return ValidationResponse{
		Message: "Validation failed",
		Errors:  validationErrors,
	}
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(fieldError validator.FieldError) string {
	field := strings.ToLower(fieldError.Field())

	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, fieldError.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, fieldError.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, fieldError.Param())
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must be a number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fieldError.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fieldError.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fieldError.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fieldError.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fieldError.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
