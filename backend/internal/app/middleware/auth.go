package middleware

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	gormaudit "github.com/potibm/kasseapparat/internal/app/store/gorm"
)

const (
	IdentityKey      = "username"
	AuthUserKey      = "auth_user"
	RemoteUserHeader = "X-Remote-User"
	DefaultUsername  = "anonymous"
)

func HandlerMiddleWare(cfg config.Config) gin.HandlerFunc {
	return ProxyAuthMiddleware(cfg)
}

func ProxyAuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Auth.Mode != "proxy" {
			c.Next()

			return
		}

		username := c.GetHeader(cfg.Auth.ProxyHeader)
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "authentication required: missing proxy header",
			})
			c.Abort()

			return
		}

		role := "user"
		if slices.Contains(cfg.Auth.ProxyAdmins, username) {
			role = "admin"
		}

		authUser := models.AuthUser{
			Username: username,
			Role:     role,
		}

		c.Set(IdentityKey, username)
		c.Set(AuthUserKey, authUser)

		ctx := gormaudit.WithUserID(c.Request.Context(), username)
		c.Request = c.Request.WithContext(ctx)

		slog.Debug("Proxy authentication successful", "username", username, "role", role)

		c.Next()
	}
}

func GetAuthUser(c *gin.Context) (*models.AuthUser, bool) {
	authUser, exists := c.Get(AuthUserKey)
	if !exists {
		return nil, false
	}

	user, ok := authUser.(models.AuthUser)
	if !ok {
		return nil, false
	}

	return &user, true
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
