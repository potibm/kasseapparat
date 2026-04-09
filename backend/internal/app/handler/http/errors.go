package http

import (
	"net/http"

	"github.com/potibm/kasseapparat/internal/app/errors"
)

func NewHTTPError(code int, message, detail string) *errors.BasicError {
	return &errors.BasicError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// predefined HTTP errors.
var (
	InvalidRequest = NewHTTPError(
		http.StatusBadRequest,
		"Invalid Request",
		"The request could not be understood by the server.",
	)
	NotFound = NewHTTPError(
		http.StatusNotFound,
		"Not Found",
		"The requested resource could not be found.",
	)
	UnableToRetrieveExecutingUser = NewHTTPError(http.StatusUnauthorized, "Unable to retrieve executing user", "")
	InternalServerError           = NewHTTPError(http.StatusInternalServerError, "Internal Server Error", "")
	Forbidden                     = NewHTTPError(
		http.StatusForbidden,
		"Forbidden",
		"You do not have permission to access this resource.",
	)
	BadRequest = NewHTTPError(
		http.StatusBadRequest,
		"Bad Request",
		"The request could not be understood by the server.",
	)
)
