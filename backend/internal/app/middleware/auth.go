package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	IdentityKey      = "username"
	RemoteUserHeader = "X-Remote-User"
	DefaultUsername  = "anonymous"
)

func HandlerMiddleWare() gin.HandlerFunc {
	return PlaceholderAuthMiddleware()
}

func PlaceholderAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader(RemoteUserHeader)
		if username == "" {
			username = DefaultUsername
			slog.Debug("No remote user header found, using default", "username", username)
		}

		c.Set(IdentityKey, username)
		c.Next()
	}
}

func Unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func GetUsername(c *gin.Context) string {
	username, exists := c.Get(IdentityKey)
	if !exists {
		return DefaultUsername
	}

	usernameStr, ok := username.(string)
	if !ok {
		return DefaultUsername
	}

	return usernameStr
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := GetUsername(c)
		if username == "" || username == DefaultUsername {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "authentication required",
			})
			c.Abort()

			return
		}

		c.Next()
	}
}
