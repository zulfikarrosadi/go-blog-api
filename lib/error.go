package lib

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorDetail struct {
	Path    []string `json:"path"`
	Value   string   `json:"value"`
	Message string   `json:"message"`
}

func ValidateError(validationError validator.ValidationErrors) []ErrorDetail {
	errorDetails := []ErrorDetail{}

	for _, fieldError := range validationError {
		switch fieldError.Tag() {
		case "required":
			errorDetail := ErrorDetail{
				Path:    []string{fieldError.Field()},
				Message: fieldError.Field() + " is required",
			}
			errorDetails = append(errorDetails, errorDetail)
		case "email":
			errorDetail := ErrorDetail{
				Path:    []string{fieldError.Field()},
				Message: "invalid email format",
			}
			errorDetails = append(errorDetails, errorDetail)
		case "eqfield":
			fmt.Println(fieldError.Field(), fieldError.StructField())
			if fieldError.Field() == "passwordConfirmation" {
				errorDetail := ErrorDetail{
					Path:    []string{"password", fieldError.Field()},
					Message: "password and password confirmation is not match",
				}
				errorDetails = append(errorDetails, errorDetail)
			}
		}
	}
	return errorDetails
}
