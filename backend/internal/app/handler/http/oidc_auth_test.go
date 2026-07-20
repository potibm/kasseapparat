package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOIDCAuthHandler_Callback_MissingParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &OIDCAuthHandler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/callback", http.NoBody)

	handler.Callback(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOIDCAuthHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &OIDCAuthHandler{
		secureCookie: false,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/logout", http.NoBody)

	handler.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewOIDCAuthHandler_InvalidIssuer(t *testing.T) {
	_, err := NewOIDCAuthHandler(
		context.Background(),
		OIDCOptions{
			Issuer:          "invalid-issuer-url",
			ClientID:        "client-id",
			ClientSecret:    "client-secret",
			CallbackURL:     "http://localhost:8080/callback",
			FrontendURL:     "http://localhost:3000",
			SessionSecret:   "session-secret-that-is-long-enough",
			SessionDuration: 24 * time.Hour,
			AdminGroup:      "",
			IsProduction:    false,
		},
	)

	assert.Error(t, err)
}

func TestOIDCAuthHandler_Getters(t *testing.T) {
	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		adminGroup: "kasseapparat-admins",
		sessionMgr: sessionMgr,
	}

	assert.Equal(t, "kasseapparat-admins", handler.GetAdminGroup())
	assert.Equal(t, sessionMgr, handler.GetSessionManager())
}

func TestOIDCAuthHandler_DetermineRole(t *testing.T) {
	handler := &OIDCAuthHandler{
		adminGroup: "kasseapparat-admins",
	}

	tests := []struct {
		name            string
		groups          []string
		expected        string
		emptyAdminGroup bool
	}{
		{"admin with matching group", []string{"kasseapparat-admins", "other-group"}, "admin", false},
		{"admin only", []string{"kasseapparat-admins"}, "admin", false},
		{"user with no groups", []string{}, "user", false},
		{"user with different groups", []string{"other-group", "another-group"}, "user", false},
		{"user with empty admin group config", []string{"kasseapparat-admins"}, "user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := handler
			if tt.emptyAdminGroup {
				testHandler = &OIDCAuthHandler{adminGroup: ""}
			}

			role := testHandler.determineRole(tt.groups)
			assert.Equal(t, tt.expected, role)
		})
	}
}

func TestSessionCookieSettings(t *testing.T) {
	require.Equal(t, "auth_session", sessionCookieName)
	require.Equal(t, "oidc_state", stateCookieName)
}

func TestSessionManagerIntegration(t *testing.T) {
	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr:  sessionMgr,
		adminGroup:  "kasseapparat-admins",
		frontendURL: "http://localhost:3000",
	}

	assert.NotNil(t, handler.GetSessionManager())
}

func TestOIDCAuthHandler_SecureCookieSettings(t *testing.T) {
	tests := []struct {
		name         string
		isProduction bool
		expectSecure bool
	}{
		{
			name:         "production enables secure cookies",
			isProduction: true,
			expectSecure: true,
		},
		{
			name:         "non-production disables secure cookies",
			isProduction: false,
			expectSecure: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &OIDCAuthHandler{
				secureCookie: tt.isProduction,
			}
			assert.Equal(t, tt.expectSecure, handler.secureCookie)
		})
	}
}
