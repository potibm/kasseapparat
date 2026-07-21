package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authMode       string
		proxyHeader    string
		proxyAdmins    []string
		headerValue    string
		expectedStatus int
		expectedUser   string
		expectedRole   string
	}{
		{
			name:           "valid user header",
			authMode:       "proxy",
			proxyHeader:    "X-Remote-User",
			proxyAdmins:    []string{},
			headerValue:    "testuser",
			expectedStatus: http.StatusOK,
			expectedUser:   "testuser",
			expectedRole:   "user",
		},
		{
			name:           "valid admin header",
			authMode:       "proxy",
			proxyHeader:    "X-Remote-User",
			proxyAdmins:    []string{"adminuser"},
			headerValue:    "adminuser",
			expectedStatus: http.StatusOK,
			expectedUser:   "adminuser",
			expectedRole:   "admin",
		},
		{
			name:           "missing header returns 401",
			authMode:       "proxy",
			proxyHeader:    "X-Remote-User",
			proxyAdmins:    []string{},
			headerValue:    "",
			expectedStatus: http.StatusUnauthorized,
			expectedUser:   "",
			expectedRole:   "",
		},
		{
			name:           "non-proxy mode passes through",
			authMode:       "oidc",
			proxyHeader:    "X-Remote-User",
			proxyAdmins:    []string{},
			headerValue:    "",
			expectedStatus: http.StatusOK,
			expectedUser:   "",
			expectedRole:   "",
		},
		{
			name:           "custom proxy header",
			authMode:       "proxy",
			proxyHeader:    "X-Custom-User",
			proxyAdmins:    []string{},
			headerValue:    "customuser",
			expectedStatus: http.StatusOK,
			expectedUser:   "customuser",
			expectedRole:   "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				Auth: config.AuthConfig{
					Mode:        tt.authMode,
					ProxyHeader: tt.proxyHeader,
					ProxyAdmins: tt.proxyAdmins,
				},
			}

			router := gin.New()
			router.Use(ProxyAuthMiddleware(cfg))
			router.GET("/test", func(c *gin.Context) {
				if tt.expectedUser != "" {
					user, ok := GetAuthUser(c)
					require.True(t, ok)
					assert.Equal(t, tt.expectedUser, user.Username)
					assert.Equal(t, tt.expectedRole, user.Role)
				}

				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
			if tt.headerValue != "" {
				req.Header.Set(tt.proxyHeader, tt.headerValue)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestProxyAuthMiddleware_SetsContextValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:        "proxy",
			ProxyHeader: "X-Remote-User",
			ProxyAdmins: []string{"admin"},
		},
	}

	router := gin.New()
	router.Use(ProxyAuthMiddleware(cfg))

	var capturedUser *models.AuthUser

	var capturedUsername string

	router.GET("/test", func(c *gin.Context) {
		user, ok := GetAuthUser(c)
		if ok {
			capturedUser = user
		}

		capturedUsername = GetUsername(c)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("X-Remote-User", "testuser")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, capturedUser)
	assert.Equal(t, "testuser", capturedUser.Username)
	assert.Equal(t, "user", capturedUser.Role)
	assert.Equal(t, "testuser", capturedUsername)
}

func TestProxyAuthMiddleware_AdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:        "proxy",
			ProxyHeader: "X-Remote-User",
			ProxyAdmins: []string{"admin1", "admin2"},
		},
	}

	router := gin.New()
	router.Use(ProxyAuthMiddleware(cfg))

	var capturedRole string

	router.GET("/test", func(c *gin.Context) {
		user, _ := GetAuthUser(c)
		capturedRole = user.Role

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("X-Remote-User", "admin1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "admin", capturedRole)
}

func TestProxyAuthMiddleware_ErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:        "proxy",
			ProxyHeader: "X-Remote-User",
		},
	}

	router := gin.New()
	router.Use(ProxyAuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, float64(http.StatusUnauthorized), response["code"])
	assert.Contains(t, response["message"], "authentication required")
}

func TestGetAuthUser_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	user, ok := GetAuthUser(c)
	assert.False(t, ok)
	assert.Nil(t, user)
}

func TestGetUsername_NotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	username := GetUsername(c)
	assert.Equal(t, DefaultUsername, username)
}

func TestHandlerMiddleWare_ProxyMode(t *testing.T) {
	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode: "proxy",
		},
	}

	middleware := HandlerMiddleWare(cfg)
	assert.NotNil(t, middleware)
}

func TestHandlerMiddleWare_OIDCMode(t *testing.T) {
	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:            "oidc",
			SessionSecret:   "test-secret-that-is-long-enough-for-testing",
			SessionDuration: 24 * time.Hour,
		},
	}

	middleware := HandlerMiddleWare(cfg)
	assert.NotNil(t, middleware)
}

func TestOIDCAuthMiddleware_MissingSessionCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:            "oidc",
			SessionSecret:   "test-secret-that-is-long-enough-for-testing",
			SessionDuration: 24 * time.Hour,
		},
	}

	router := gin.New()
	router.Use(OIDCAuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["message"], "missing session cookie")
}

func TestOIDCAuthMiddleware_InvalidSessionCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:            "oidc",
			SessionSecret:   "test-secret-that-is-long-enough-for-testing",
			SessionDuration: 24 * time.Hour,
		},
	}

	router := gin.New()
	router.Use(OIDCAuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
		Name:     "auth_session",
		Value:    "invalid-session-data",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["message"], "invalid session")
}

func TestOIDCAuthMiddleware_NonOIDCMode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:        "proxy",
			ProxyHeader: "X-Remote-User",
		},
	}

	router := gin.New()
	router.Use(OIDCAuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOIDCAuthMiddleware_ValidSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)
	sessionData := session.SessionData{
		Username:  "testuser",
		Role:      "admin",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	encodedSession, err := sessionMgr.EncodeSession(sessionData)
	require.NoError(t, err)

	cfg := config.Config{
		Auth: config.AuthConfig{
			Mode:            "oidc",
			SessionSecret:   "test-secret-that-is-long-enough-for-testing",
			SessionDuration: 24 * time.Hour,
		},
	}

	router := gin.New()
	router.Use(OIDCAuthMiddleware(cfg))

	var capturedUser *models.AuthUser

	var capturedUsername string

	router.GET("/test", func(c *gin.Context) {
		user, ok := GetAuthUser(c)
		if ok {
			capturedUser = user
		}

		capturedUsername = GetUsername(c)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
		Name:     "auth_session",
		Value:    encodedSession,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, capturedUser)
	assert.Equal(t, "testuser", capturedUser.Username)
	assert.Equal(t, "admin", capturedUser.Role)
	assert.Equal(t, "testuser", capturedUsername)
}

func TestGetAuthUser_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	c.Set(AuthUserKey, "not-an-auth-user")

	user, ok := GetAuthUser(c)
	assert.False(t, ok)
	assert.Nil(t, user)
}

func TestGetUsername_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	c.Set(IdentityKey, 12345)

	username := GetUsername(c)
	assert.Equal(t, DefaultUsername, username)
}
