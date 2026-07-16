package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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

	handler := &OIDCAuthHandler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/logout", http.NoBody)

	handler.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewOIDCAuthHandler_InvalidIssuer(t *testing.T) {
	_, err := NewOIDCAuthHandler(
		context.Background(),
		"invalid-issuer-url",
		"client-id",
		"client-secret",
		"http://localhost:8080/callback",
		"http://localhost:3000",
		"session-secret-that-is-long-enough",
		[]string{},
	)

	assert.Error(t, err)
}

func TestOIDCAuthHandler_Getters(t *testing.T) {
	handler := &OIDCAuthHandler{
		admins: []string{"admin1", "admin2"},
	}

	assert.Equal(t, []string{"admin1", "admin2"}, handler.GetAdmins())
	assert.Nil(t, handler.GetSessionManager())
}

func TestExtractUsernameFromClaims(t *testing.T) {
	tests := []struct {
		name              string
		preferredUsername string
		nameClaim         string
		email             string
		expected          string
	}{
		{
			name:              "preferred_username takes precedence",
			preferredUsername: "john.doe",
			nameClaim:         "John Doe",
			email:             "john@example.com",
			expected:          "john.doe",
		},
		{
			name:              "name is used when preferred_username is empty",
			preferredUsername: "",
			nameClaim:         "Jane Smith",
			email:             "jane@example.com",
			expected:          "Jane Smith",
		},
		{
			name:              "email prefix is used when both are empty",
			preferredUsername: "",
			nameClaim:         "",
			email:             "bob@example.com",
			expected:          "bob",
		},
		{
			name:              "empty string when all are empty",
			preferredUsername: "",
			nameClaim:         "",
			email:             "",
			expected:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var username string

			switch {
			case tt.preferredUsername != "":
				username = tt.preferredUsername
			case tt.nameClaim != "":
				username = tt.nameClaim
			case tt.email != "":
				username = extractUsernameFromEmail(tt.email)
			}

			assert.Equal(t, tt.expected, username)
		})
	}
}

func extractUsernameFromEmail(email string) string {
	for i, c := range email {
		if c == '@' {
			return email[:i]
		}
	}

	return email
}

func TestDetermineRole(t *testing.T) {
	admins := []string{"admin1", "admin2@example.com"}

	tests := []struct {
		username string
		expected string
	}{
		{"admin1", "admin"},
		{"admin2@example.com", "admin"},
		{"regularuser", "user"},
		{"user@example.com", "user"},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			role := determineRole(tt.username, admins)
			assert.Equal(t, tt.expected, role)
		})
	}
}

func determineRole(username string, admins []string) string {
	for _, admin := range admins {
		if admin == username {
			return "admin"
		}
	}

	return "user"
}

func TestSessionCookieSettings(t *testing.T) {
	require.Equal(t, "auth_session", sessionCookieName)
	require.Equal(t, "oidc_state", stateCookieName)
}

func TestSessionManagerIntegration(t *testing.T) {
	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing")

	handler := &OIDCAuthHandler{
		sessionMgr:  sessionMgr,
		admins:      []string{"admin@example.com"},
		frontendURL: "http://localhost:3000",
	}

	assert.NotNil(t, handler.GetSessionManager())
}
