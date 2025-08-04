package middleware

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type HttpError interface {
	StatusCode() int
	Error() string
	Details() string
	Cause() error
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0].Err
		hub := sentrygin.GetHubFromContext(c)

		// capture the error in Sentry
		captureError(hub, err)

		// return error response to the client
		writeErrorResponse(c, err)
	}
}

func captureError(hub *sentry.Hub, err error) {
	if hub == nil || err == nil {
		return
	}

	// unwrap known HttpError
	if httpErr, ok := err.(HttpError); ok {
		if cause := httpErr.Cause(); cause != nil {
			hub.CaptureException(cause)

			return
		}
	}

	hub.CaptureException(err)
}

func writeErrorResponse(c *gin.Context, err error) {
	if httpErr, ok := err.(HttpError); ok {
		response := gin.H{
			"error": httpErr.Error(),
		}
		if detail := httpErr.Details(); detail != "" {
			response["details"] = detail
		}

		c.JSON(httpErr.StatusCode(), response)

		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal Server Error",
	})
}
