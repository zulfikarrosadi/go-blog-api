package lib

import (
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
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

func ValidateErrorV2(action string, err error) {
	fmt.Printf("error in %v: %v", action, err)
	// check for database and query error
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		switch driverErr.Number {
		case mysqlerr.ER_SYNTAX_ERROR:
			ErrorLog(action, "sql syntax error", err)
		case mysqlerr.ER_ACCESS_DENIED_ERROR:
			ErrorLog(action, "database username or password is incorrect", err)
		case mysqlerr.ER_NO_SUCH_TABLE:
			ErrorLog(action, "table not exists", err)
		case mysqlerr.ER_DUP_ENTRY:
			ErrorLog(action, "username already exists", err)
		default:
			ErrorLog(action, err.Error(), err)
		}
	}

	// check for database auth error
	if errors.Is(err, sql.ErrNoRows) {
		ErrorLog(action, "username or password is incorrect", err)
	}

	// check for database connection error
	// eg: port error, protocol error, etc
	if newErr := err.(*net.OpError); newErr != nil {
		ErrorLog(action, newErr.Error(), newErr)
	}
}
