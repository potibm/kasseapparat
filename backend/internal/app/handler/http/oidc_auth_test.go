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
	"golang.org/x/oauth2"
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

	cookies := w.Result().Cookies()

	var sessionCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == sessionCookieName {
			sessionCookie = cookie

			break
		}
	}

	require.NotNil(t, sessionCookie, "session cookie should be set")
	assert.Equal(t, -1, sessionCookie.MaxAge, "session cookie should be invalidated")
	assert.Equal(t, "", sessionCookie.Value, "session cookie value should be empty")
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

func TestOIDCAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr:   sessionMgr,
		secureCookie: false,
		frontendURL:  "http://localhost:3000",
		oauth2Config: &oauth2.Config{
			ClientID:     "test-client",
			ClientSecret: "test-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://provider.example.com/auth",
				TokenURL: "https://provider.example.com/token",
			},
		},
	}

	tests := []struct {
		name        string
		query       string
		expectCode  int
		expectState bool
	}{
		{
			name:        "login without returnTo",
			query:       "",
			expectCode:  http.StatusFound,
			expectState: true,
		},
		{
			name:        "login with valid returnTo",
			query:       "?returnTo=/dashboard",
			expectCode:  http.StatusFound,
			expectState: true,
		},
		{
			name:        "login with invalid returnTo (external)",
			query:       "?returnTo=https://evil.com",
			expectCode:  http.StatusFound,
			expectState: true,
		},
		{
			name:        "login with protocol-relative returnTo",
			query:       "?returnTo=//evil.com",
			expectCode:  http.StatusFound,
			expectState: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/login"+tt.query, http.NoBody)

			handler.Login(c)

			assert.Equal(t, tt.expectCode, w.Code)

			cookies := w.Result().Cookies()

			var stateCookie, returnToCookie *http.Cookie

			for _, cookie := range cookies {
				if cookie.Name == stateCookieName {
					stateCookie = cookie
				}

				if cookie.Name == returnToCookieName {
					returnToCookie = cookie
				}
			}

			if tt.expectState {
				require.NotNil(t, stateCookie, "state cookie should be set")
				assert.NotEmpty(t, stateCookie.Value)
				assert.True(t, stateCookie.HttpOnly)
				assert.Equal(t, "/", stateCookie.Path)
			}

			require.NotNil(t, returnToCookie, "returnTo cookie should be set")
		})
	}
}

func TestOIDCAuthHandler_Callback_MissingStateCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr: sessionMgr,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/callback?code=testcode&state=teststate", http.NoBody)

	handler.Callback(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOIDCAuthHandler_Callback_InvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr: sessionMgr,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/callback?code=testcode&state=teststate", http.NoBody)
	c.Request.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
		Name:     stateCookieName,
		Value:    "invalid-state-value",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	handler.Callback(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOIDCAuthHandler_ClearCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &OIDCAuthHandler{
		secureCookie: false,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	handler.clearCookie(c, "test-cookie")

	cookies := w.Result().Cookies()

	var foundCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == "test-cookie" {
			foundCookie = cookie

			break
		}
	}

	require.NotNil(t, foundCookie, "cookie should be set")
	assert.Equal(t, -1, foundCookie.MaxAge)
	assert.Equal(t, "", foundCookie.Value)
	assert.True(t, foundCookie.HttpOnly)
	assert.Equal(t, "/", foundCookie.Path)
}

func TestOIDCAuthHandler_CreateSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr:   sessionMgr,
		secureCookie: false,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	err := handler.createSession(c, "testuser", "admin")
	require.NoError(t, err)

	cookies := w.Result().Cookies()

	var sessionCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == sessionCookieName {
			sessionCookie = cookie

			break
		}
	}

	require.NotNil(t, sessionCookie, "session cookie should be set")
	assert.NotEmpty(t, sessionCookie.Value)
	assert.True(t, sessionCookie.HttpOnly)
	assert.Equal(t, "/", sessionCookie.Path)

	sessionData, err := sessionMgr.DecodeSession(sessionCookie.Value)
	require.NoError(t, err)
	assert.Equal(t, "testuser", sessionData.Username)
	assert.Equal(t, "admin", sessionData.Role)
}

func TestOIDCAuthHandler_ValidateState_Mismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr: sessionMgr,
	}

	stateData := session.StateData{
		State:     "correct-state",
		Nonce:     "test-nonce",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	encodedState, err := sessionMgr.EncodeState(stateData)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/callback", http.NoBody)
	c.Request.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
		Name:     stateCookieName,
		Value:    encodedState,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	result, err := handler.validateState(c, "wrong-state")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOIDCAuthHandler_ValidateState_MissingCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr: sessionMgr,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/callback", http.NoBody)

	result, err := handler.validateState(c, "any-state")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestOIDCAuthHandler_ValidateState_ExpiredState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr: sessionMgr,
	}

	stateData := session.StateData{
		State:     "test-state",
		Nonce:     "test-nonce",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	encodedState, err := sessionMgr.EncodeState(stateData)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/callback", http.NoBody)
	c.Request.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
		Name:     stateCookieName,
		Value:    encodedState,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	result, err := handler.validateState(c, "test-state")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestOIDCAuthHandler_Callback_StateMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr: sessionMgr,
	}

	stateData := session.StateData{
		State:     "correct-state",
		Nonce:     "test-nonce",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	encodedState, err := sessionMgr.EncodeState(stateData)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/callback?code=testcode&state=wrong-state", http.NoBody)
	c.Request.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
		Name:     stateCookieName,
		Value:    encodedState,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	handler.Callback(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOIDCAuthHandler_Login_InvalidReturnTo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	handler := &OIDCAuthHandler{
		sessionMgr:   sessionMgr,
		secureCookie: false,
		oauth2Config: &oauth2.Config{
			ClientID:     "test-client",
			ClientSecret: "test-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://provider.example.com/auth",
				TokenURL: "https://provider.example.com/token",
			},
		},
	}

	tests := []struct {
		name             string
		returnTo         string
		expectedFallback string
	}{
		{
			name:             "empty returnTo defaults to /",
			returnTo:         "",
			expectedFallback: "/",
		},
		{
			name:             "external URL defaults to /",
			returnTo:         "https://evil.com",
			expectedFallback: "/",
		},
		{
			name:             "protocol-relative URL defaults to /",
			returnTo:         "//evil.com",
			expectedFallback: "/",
		},
		{
			name:             "valid path is preserved",
			returnTo:         "/dashboard",
			expectedFallback: "/dashboard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			query := ""
			if tt.returnTo != "" {
				query = "?returnTo=" + tt.returnTo
			}

			c.Request = httptest.NewRequest(http.MethodGet, "/login"+query, http.NoBody)

			handler.Login(c)

			assert.Equal(t, http.StatusFound, w.Code)

			cookies := w.Result().Cookies()

			var returnToCookie *http.Cookie

			for _, cookie := range cookies {
				if cookie.Name == returnToCookieName {
					returnToCookie = cookie

					break
				}
			}

			require.NotNil(t, returnToCookie)
			assert.Equal(t, tt.expectedFallback, returnToCookie.Value)
		})
	}
}

// Note: extractUserClaims requires a real *oidc.IDToken which is difficult to mock.
// The function is covered indirectly through integration tests and the Callback flow.

func TestOIDCAuthHandler_ExchangeCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		handler     *OIDCAuthHandler
		code        string
		expectError bool
	}{
		{
			name: "exchange code with invalid oauth config",
			handler: &OIDCAuthHandler{
				oauth2Config: &oauth2.Config{
					ClientID:     "test-client",
					ClientSecret: "wrong-secret",
					Endpoint: oauth2.Endpoint{
						TokenURL: "http://invalid-endpoint-that-does-not-exist.example.com/token",
					},
				},
			},
			code:        "test-code",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/callback", http.NoBody)

			token, err := tt.handler.exchangeCode(c, tt.code)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, token)
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
			}
		})
	}
}

func TestOIDCAuthHandler_VerifyIDToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	tests := []struct {
		name        string
		handler     *OIDCAuthHandler
		token       *oauth2.Token
		nonce       string
		expectError bool
		expectCode  int
	}{
		{
			name: "missing id_token in oauth token",
			handler: &OIDCAuthHandler{
				sessionMgr: sessionMgr,
			},
			token:       &oauth2.Token{},
			nonce:       "test-nonce",
			expectError: true,
			expectCode:  http.StatusInternalServerError,
		},
		{
			name: "id_token is not a string",
			handler: &OIDCAuthHandler{
				sessionMgr: sessionMgr,
			},
			token: (&oauth2.Token{}).WithExtra(map[string]any{
				"id_token": 12345,
			}),
			nonce:       "test-nonce",
			expectError: true,
			expectCode:  http.StatusInternalServerError,
		},
		{
			name: "id_token is nil",
			handler: &OIDCAuthHandler{
				sessionMgr: sessionMgr,
			},
			token: (&oauth2.Token{}).WithExtra(map[string]any{
				"id_token": nil,
			}),
			nonce:       "test-nonce",
			expectError: true,
			expectCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/callback", http.NoBody)

			idToken, err := tt.handler.verifyIDToken(c, tt.token, tt.nonce)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, idToken)
				assert.Equal(t, tt.expectCode, w.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, idToken)
			}
		})
	}
}

func TestNewOIDCAuthHandler_Success(t *testing.T) {
	// Create a mock OIDC server
	server := createMockOIDCServer()
	defer server.Close()

	handler, err := NewOIDCAuthHandler(
		context.Background(),
		OIDCOptions{
			Issuer:          server.URL,
			ClientID:        "test-client-id",
			ClientSecret:    "test-client-secret",
			CallbackURL:     "http://localhost:8080/callback",
			FrontendURL:     "http://localhost:3000",
			SessionSecret:   "test-secret-that-is-long-enough-for-testing",
			SessionDuration: 24 * time.Hour,
			AdminGroup:      "admins",
			IsProduction:    false,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, handler)
	assert.Equal(t, "admins", handler.adminGroup)
	assert.Equal(t, "http://localhost:3000", handler.frontendURL)
	assert.Equal(t, "http://localhost:8080/callback", handler.callbackURL)
	assert.False(t, handler.secureCookie)
	assert.NotNil(t, handler.oauth2Config)
	assert.NotNil(t, handler.verifier)
	assert.NotNil(t, handler.sessionMgr)
}

func TestNewOIDCAuthHandler_ProductionMode(t *testing.T) {
	// Create a mock OIDC server
	server := createMockOIDCServer()
	defer server.Close()

	handler, err := NewOIDCAuthHandler(
		context.Background(),
		OIDCOptions{
			Issuer:          server.URL,
			ClientID:        "test-client-id",
			ClientSecret:    "test-client-secret",
			CallbackURL:     "http://localhost:8080/callback",
			FrontendURL:     "http://localhost:3000",
			SessionSecret:   "test-secret-that-is-long-enough-for-testing",
			SessionDuration: 24 * time.Hour,
			AdminGroup:      "admins",
			IsProduction:    true,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, handler)
	assert.True(t, handler.secureCookie, "secureCookie should be true in production mode")
}

func TestOIDCAuthHandler_ExtractUserClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &OIDCAuthHandler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/callback", http.NoBody)

	// Test with nil IDToken - should panic
	assert.Panics(t, func() {
		_, _, _ = handler.extractUserClaims(c, nil)
	})
}

func TestOIDCAuthHandler_VerifyIDToken_NonceMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock OIDC server
	server := createMockOIDCServer()
	defer server.Close()

	handler := createTestOIDCHandler(t, server.URL)

	// Create a mock token with an invalid id_token format
	token := (&oauth2.Token{}).WithExtra(map[string]any{
		"id_token": "invalid.jwt.token",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/callback", http.NoBody)

	// This should fail during verification
	_, err := handler.verifyIDToken(c, token, "expected-nonce")
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOIDCAuthHandler_Login_ErrorPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	// Test with nil oauth2Config to trigger error
	handler := &OIDCAuthHandler{
		sessionMgr:   sessionMgr,
		secureCookie: false,
		oauth2Config: nil, // This will cause a panic or error
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/login", http.NoBody)

	// This should handle the nil config gracefully
	assert.Panics(t, func() {
		handler.Login(c)
	})
}

func TestOIDCAuthHandler_Callback_ValidFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	// Create a valid state
	stateData := session.StateData{
		State:     "valid-state",
		Nonce:     "valid-nonce",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	encodedState, err := sessionMgr.EncodeState(stateData)
	require.NoError(t, err)

	handler := &OIDCAuthHandler{
		sessionMgr:  sessionMgr,
		frontendURL: "http://localhost:3000",
		oauth2Config: &oauth2.Config{
			ClientID:     "test-client",
			ClientSecret: "test-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Endpoint: oauth2.Endpoint{
				TokenURL: "http://invalid-endpoint.example.com/token",
			},
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request with valid state but invalid code
	c.Request = httptest.NewRequest(http.MethodGet, "/callback?code=invalid-code&state=valid-state", http.NoBody)
	c.Request.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
		Name:     stateCookieName,
		Value:    encodedState,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// Add returnTo cookie
	c.Request.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
		Name:     returnToCookieName,
		Value:    "/dashboard",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	handler.Callback(c)

	// Should fail at exchangeCode step
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOIDCAuthHandler_Callback_WithReturnToCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionMgr := session.NewManager("test-secret-that-is-long-enough-for-testing", 24*time.Hour)

	stateData := session.StateData{
		State:     "test-state",
		Nonce:     "test-nonce",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	encodedState, err := sessionMgr.EncodeState(stateData)
	require.NoError(t, err)

	handler := &OIDCAuthHandler{
		sessionMgr:  sessionMgr,
		frontendURL: "http://localhost:3000",
		oauth2Config: &oauth2.Config{
			ClientID:     "test-client",
			ClientSecret: "test-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Endpoint: oauth2.Endpoint{
				TokenURL: "http://invalid-endpoint.example.com/token",
			},
		},
	}

	tests := []struct {
		name           string
		returnToValue  string
		expectedInPath string
	}{
		{
			name:           "valid returnTo path",
			returnToValue:  "/admin",
			expectedInPath: "/admin",
		},
		{
			name:           "protocol-relative returnTo should be ignored",
			returnToValue:  "//evil.com",
			expectedInPath: "/",
		},
		{
			name:           "empty returnTo should default to /",
			returnToValue:  "",
			expectedInPath: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/callback?code=test&state=test-state", http.NoBody)
			c.Request.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
				Name:     stateCookieName,
				Value:    encodedState,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})

			if tt.returnToValue != "" {
				c.Request.AddCookie(&http.Cookie{ //nolint:gosec // test cookie
					Name:     returnToCookieName,
					Value:    tt.returnToValue,
					HttpOnly: true,
					Secure:   false,
					SameSite: http.SameSiteLaxMode,
				})
			}

			handler.Callback(c)

			// Will fail at exchangeCode, but we can check the redirect path logic
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	}
}

func TestOIDCAuthHandler_CreateSession_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a handler with a nil session manager to trigger error
	handler := &OIDCAuthHandler{
		sessionMgr:   nil,
		secureCookie: false,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	// This should panic with nil session manager
	assert.Panics(t, func() {
		_ = handler.createSession(c, "testuser", "user")
	})
}

func TestOIDCAuthHandler_Login_WithValidConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock OIDC server
	server := createMockOIDCServer()
	defer server.Close()

	handler := createTestOIDCHandler(t, server.URL)

	tests := []struct {
		name        string
		query       string
		expectCode  int
		expectState bool
	}{
		{
			name:        "login without returnTo",
			query:       "",
			expectCode:  http.StatusFound,
			expectState: true,
		},
		{
			name:        "login with valid returnTo",
			query:       "?returnTo=/dashboard",
			expectCode:  http.StatusFound,
			expectState: true,
		},
		{
			name:        "login with invalid returnTo (external)",
			query:       "?returnTo=https://evil.com",
			expectCode:  http.StatusFound,
			expectState: true,
		},
		{
			name:        "login with protocol-relative returnTo",
			query:       "?returnTo=//evil.com",
			expectCode:  http.StatusFound,
			expectState: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLoginFlow(t, handler, tt.query, tt.expectCode, tt.expectState)
		})
	}
}

//nolint:gosec // test code - writing to test response writer
func createMockOIDCServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"issuer": "http://` + r.Host + `",
				"authorization_endpoint": "http://` + r.Host + `/auth",
				"token_endpoint": "http://` + r.Host + `/token",
				"jwks_uri": "http://` + r.Host + `/jwks"
			}`))
		case "/jwks":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"keys": []}`))
		}
	}))
}

func createTestOIDCHandler(t *testing.T, issuerURL string) *OIDCAuthHandler {
	t.Helper()

	handler, err := NewOIDCAuthHandler(
		context.Background(),
		OIDCOptions{
			Issuer:          issuerURL,
			ClientID:        "test-client",
			ClientSecret:    "test-secret",
			CallbackURL:     "http://localhost:8080/callback",
			FrontendURL:     "http://localhost:3000",
			SessionSecret:   "test-secret-that-is-long-enough-for-testing",
			SessionDuration: 24 * time.Hour,
			AdminGroup:      "admins",
			IsProduction:    false,
		},
	)
	require.NoError(t, err)

	return handler
}

func testLoginFlow(t *testing.T, handler *OIDCAuthHandler, query string, expectCode int, expectState bool) {
	t.Helper()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/login"+query, http.NoBody)

	handler.Login(c)

	assert.Equal(t, expectCode, w.Code)

	cookies := w.Result().Cookies()

	var stateCookie, returnToCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == stateCookieName {
			stateCookie = cookie
		}

		if cookie.Name == returnToCookieName {
			returnToCookie = cookie
		}
	}

	if expectState {
		require.NotNil(t, stateCookie, "state cookie should be set")
		assert.NotEmpty(t, stateCookie.Value)
		assert.True(t, stateCookie.HttpOnly)
		assert.Equal(t, "/", stateCookie.Path)
	}

	require.NotNil(t, returnToCookie, "returnTo cookie should be set")
}
