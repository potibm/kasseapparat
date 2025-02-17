package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpError interface {
	StatusCode() int
	Error() string
	Details() string
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err

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
	}
}
