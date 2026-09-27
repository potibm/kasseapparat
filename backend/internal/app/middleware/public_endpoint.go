package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
)

// bearerScheme is the Authorization scheme the public endpoints accept.
const bearerScheme = "Bearer"

// PublicEndpointAuth guards the endpoints that must stay reachable without a login,
// such as the purchase statistics display. The caller proves it may read the data by
// presenting the configured token as `Authorization: Bearer <token>`.
//
// This is deliberately not a login: it authorises one endpoint, carries no identity
// and grants nothing else. The token is a shared secret, so anything that leaks it
// (a config file, a log line, a browser that stored it) exposes the endpoint.
func PublicEndpointAuth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !hasValidPublicEndpointToken(c, cfg.Auth.PublicEndpointToken) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "authentication required: missing or invalid bearer token",
			})

			return
		}

		c.Next()
	}
}

func hasValidPublicEndpointToken(c *gin.Context, expected string) bool {
	// Fail closed when nothing is configured. Without this an empty
	// `Authorization: Bearer ` would present an empty string and compare equal to an
	// empty expected token, leaving the endpoint open to anyone.
	if expected == "" {
		return false
	}

	scheme, presented, found := strings.Cut(c.GetHeader("Authorization"), " ")
	if !found || !strings.EqualFold(scheme, bearerScheme) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(strings.TrimSpace(presented)), []byte(expected)) == 1
}
