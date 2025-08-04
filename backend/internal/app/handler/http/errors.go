package http

import (
	"net/http"

	"github.com/potibm/kasseapparat/internal/app/errors"
)

func NewHttpError(code int, message string, detail string) *errors.BasicError {
	return &errors.BasicError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// predefined HTTP errors.
var (
	InvalidRequest                = NewHttpError(http.StatusBadRequest, "Invalid Request", "The request could not be understood by the server.")
	NotFound                      = NewHttpError(http.StatusNotFound, "Not Found", "The requested resource could not be found.")
	UnableToRetrieveExecutingUser = NewHttpError(http.StatusUnauthorized, "Unable to retrieve executing user", "")
	InternalServerError           = NewHttpError(http.StatusInternalServerError, "Internal Server Error", "")
	Forbidden                     = NewHttpError(http.StatusForbidden, "Forbidden", "You do not have permission to access this resource.")
	BadRequest                    = NewHttpError(http.StatusBadRequest, "Bad Request", "The request could not be understood by the server.")
)
