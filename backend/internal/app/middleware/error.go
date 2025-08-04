package middleware

import (
	"net/http"

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

		// we have a detailed error (HttpError aka errors.BasicError)
		if httpErr, ok := err.(HttpError); ok {
			if hub != nil {
				if cause := httpErr.Cause(); cause != nil {
					hub.CaptureException(cause)
				} else {
					hub.CaptureException(httpErr)
				}
			}

			response := gin.H{
				"error": httpErr.Error(),
			}
			if detail := httpErr.Details(); detail != "" {
				response["details"] = detail
			}

			c.JSON(httpErr.StatusCode(), response)

			return
		}

		// case of a generic error
		if hub != nil {
			hub.CaptureException(err)
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	}
}
