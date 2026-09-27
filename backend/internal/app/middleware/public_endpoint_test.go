package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func TestPublicEndpointAuth(t *testing.T) {
	const token = "s3cr3t-public-token"

	tests := []struct {
		name       string
		cfgToken   string
		header     string
		wantStatus int
	}{
		{name: "no header", cfgToken: token, wantStatus: http.StatusUnauthorized},
		{name: "empty header", cfgToken: token, header: "", wantStatus: http.StatusUnauthorized},
		{
			name:       "wrong scheme",
			cfgToken:   token,
			header:     "Basic " + token,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "scheme without token",
			cfgToken:   token,
			header:     "Bearer",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong token",
			cfgToken:   token,
			header:     "Bearer nope",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "token is a prefix of the configured one",
			cfgToken:   token,
			header:     "Bearer s3cr3t-public",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "valid token",
			cfgToken:   token,
			header:     "Bearer " + token,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "scheme is case insensitive",
			cfgToken:   token,
			header:     "bearer " + token,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "surrounding whitespace is tolerated",
			cfgToken:   token,
			header:     "Bearer   " + token + "  ",
			wantStatus: http.StatusNoContent,
		},
		// An empty configured token must never match, otherwise an empty
		// "Bearer " would open the endpoint to everyone.
		{
			name:       "no token configured and none presented",
			cfgToken:   "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "no token configured and empty one presented",
			cfgToken:   "",
			header:     "Bearer ",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "no token configured and one presented",
			cfgToken:   "",
			header:     "Bearer anything",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			router := gin.New()
			router.Use(ErrorHandlingMiddleware())
			router.GET("/stats", PublicEndpointAuth(config.Config{
				Auth: config.AuthConfig{PublicEndpointToken: tc.cfgToken},
			}), func(c *gin.Context) {
				c.Status(http.StatusNoContent)
			})

			request := httptest.NewRequest(http.MethodGet, "/stats", http.NoBody)
			if tc.header != "" {
				request.Header.Set("Authorization", tc.header)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			assert.Equal(t, tc.wantStatus, recorder.Code)
		})
	}
}
