package validation

import (
	"log"

	"github.com/go-playground/validator/v10"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

var validate = validator.New()

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func ValidationRequest(payload any) []string {
	err := validate.Struct(payload)
	if err == nil {
		return nil
	}
	validationErrors := err.(validator.ValidationErrors)
	if len(validationErrors) > 0 {
		var errorMessages []string
		for _, fieldError := range validationErrors {
			errorMessages = append(errorMessages, fieldError.Error())
		}
		log.Println("Validation errors:", errorMessages)

		return errorMessages
	}
	return nil
}
