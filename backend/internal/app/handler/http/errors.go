package http

import (
	"net/http"

	"github.com/potibm/kasseapparat/internal/app/middleware"
)

type BasicError struct {
	Code    int
	Message string
	Detail  string
}

func (e *BasicError) StatusCode() int {
	return e.Code
}

func (e *BasicError) Error() string {
	return e.Message
}

func (e *BasicError) Details() string {
	return e.Detail
}

func (e *BasicError) SetDetails(details string) *BasicError {
	newError := &BasicError{
		Code:    e.Code,
		Message: e.Message,
		Detail:  details,
	}

	return newError
}

func NewHttpError(code int, message string, detail string) middleware.HttpError {
	return &BasicError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

func ExtendHttpErrorWithDetails(httpError middleware.HttpError, message string) middleware.HttpError {
	return &BasicError{
		Code:    httpError.StatusCode(),
		Message: httpError.Error(),
		Detail:  message,
	}
}

// Vordefinierte Fehler für häufige Szenarien.
var (
	InvalidRequest                = NewHttpError(http.StatusBadRequest, "Invalid Request", "The request could not be understood by the server.")
	NotFound                      = NewHttpError(http.StatusNotFound, "Not Found", "The requested resource could not be found.")
	UnableToRetrieveExecutingUser = NewHttpError(http.StatusUnauthorized, "Unable to retrieve executing user", "")
	InternalServerError           = NewHttpError(http.StatusInternalServerError, "Internal Server Error", "")
	Forbidden                     = NewHttpError(http.StatusForbidden, "Forbidden", "You do not have permission to access this resource.")
	BadRequest                    = NewHttpError(http.StatusBadRequest, "Bad Request", "The request could not be understood by the server.")
)
