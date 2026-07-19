// Package http provides HTTP handlers for the Kasseapparat application,
// including OIDC authentication for the BFF (Backend-for-Frontend) pattern.
//
// # OIDC Authentication Flow
//
// The OIDC handler implements a secure authentication flow using the Authorization Code Grant:
//
//  1. Login: Generates cryptographically secure state and nonce values, stores them in an
//     encrypted HttpOnly cookie, and redirects to the Identity Provider (IdP).
//
//  2. Callback: Validates the state parameter (CSRF protection), exchanges the authorization
//     code for tokens, verifies the ID token and nonce (replay attack protection), extracts
//     user information from claims, creates an encrypted session cookie, and redirects to frontend.
//
//  3. Logout: Clears the session cookie to terminate the user's session.
//
// # Security Considerations
//
//   - CSRF Protection: The state parameter is generated using crypto/rand and validated on callback
//   - Replay Attack Prevention: The nonce is included in the ID token and verified
//   - Secure Cookies: Session and state cookies are HttpOnly (XSS protection), use SameSite=Lax,
//     and Secure flag in production (HTTPS only)
//   - Encrypted Sessions: Session data is encrypted using gorilla/securecookie with derived keys
//   - Configurable Duration: Session lifetime is configurable via auth.session_duration
//   - Admin Detection: Users in the configured admin list receive "admin" role, others get "user"
package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/potibm/kasseapparat/internal/app/session"
	"golang.org/x/oauth2"
)

const (
	stateCookieName    = "oidc_state"
	returnToCookieName = "oidc_return_to"
	sessionCookieName  = "auth_session"
)

type OIDCOptions struct {
	Issuer          string
	ClientID        string
	ClientSecret    string
	CallbackURL     string
	FrontendURL     string
	SessionSecret   string
	SessionDuration time.Duration
	Admins          []string
	IsProduction    bool
}

type OIDCAuthHandler struct {
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	sessionMgr   *session.Manager
	admins       []string
	frontendURL  string
	callbackURL  string
	secureCookie bool
	log          slog.Logger
}

// NewOIDCAuthHandler initializes the OIDC authentication handler with the provided options.
//
// It discovers the OIDC provider configuration from the issuer URL, sets up OAuth2 config,
// creates a session manager with encrypted cookies, and configures secure cookie settings
// based on the environment (production vs development).
//
// Security considerations:
//   - In production, cookies use the Secure flag (HTTPS only)
//   - In non-production environments, a warning is logged about insecure cookies
//   - Session secret must be at least 32 characters (validated in config)
//   - Session duration is configurable via opts.SessionDuration
func NewOIDCAuthHandler(
	ctx context.Context,
	opts OIDCOptions,
) (*OIDCAuthHandler, error) {
	provider, err := oidc.NewProvider(ctx, opts.Issuer)
	if err != nil {
		return nil, err
	}

	if !opts.IsProduction {
		slog.Warn(
			"SECURE COOKIES ARE DISABLED! This is highly insecure and must NOT be used in production environments.",
		)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURL:  opts.CallbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: opts.ClientID})
	sessionMgr := session.NewManager(opts.SessionSecret, opts.SessionDuration)

	return &OIDCAuthHandler{
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		sessionMgr:   sessionMgr,
		admins:       opts.Admins,
		frontendURL:  opts.FrontendURL,
		callbackURL:  opts.CallbackURL,
		secureCookie: opts.IsProduction,
	}, nil
}

// Login initiates the OIDC authentication flow by generating a secure state
// and nonce, storing them in an encrypted cookie, and redirecting to the IdP.
//
// Security measures:
//   - State: Cryptographically random 32-byte value prevents CSRF attacks
//   - Nonce: Cryptographically random 32-byte value prevents replay attacks
//   - Cookie: HttpOnly, Secure (in production), SameSite=Lax, encrypted with securecookie
//   - Expiration: State expires after 10 minutes (configurable via session.StateDuration)
//
// The state and nonce are stored together in an encrypted cookie to maintain
// server-side state without requiring server-side storage.
func (h *OIDCAuthHandler) Login(c *gin.Context) {
	state, err := session.GenerateRandomString(session.RandomStringLength)
	if err != nil {
		slog.Error("Failed to generate state", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})

		return
	}

	nonce, err := session.GenerateRandomString(session.RandomStringLength)
	if err != nil {
		slog.Error("Failed to generate nonce", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate nonce"})

		return
	}

	stateData := session.StateData{
		State:     state,
		Nonce:     nonce,
		ExpiresAt: time.Now().Add(session.StateDuration),
	}

	encodedState, err := h.sessionMgr.EncodeState(stateData)
	if err != nil {
		slog.Error("Failed to encode state", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode state"})

		return
	}

	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // Secure is set dynamically via h.secureCookie
		Name:     stateCookieName,
		Value:    encodedState,
		MaxAge:   int(session.StateDuration.Seconds()),
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	returnTo := c.Query("returnTo")

	if returnTo == "" || len(returnTo) > 2000 || !strings.HasPrefix(returnTo, "/") ||
		strings.HasPrefix(returnTo, "//") {
		returnTo = "/" // default fallback
	}

	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // Secure is set dynamically via h.secureCookie
		Name:     returnToCookieName,
		Value:    returnTo,
		MaxAge:   int(session.StateDuration.Seconds()),
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	authURL := h.oauth2Config.AuthCodeURL(state, oauth2.SetAuthURLParam("nonce", nonce))

	slog.Debug("Redirecting to OIDC provider", "url", authURL)

	c.Redirect(http.StatusFound, authURL)
}

// Callback handles the OIDC callback from the identity provider.
//
// Flow:
//  1. Extract authorization code and state from query parameters
//  2. Validate state parameter against encrypted cookie (CSRF protection)
//  3. Exchange authorization code for tokens with the IdP
//  4. Verify ID token signature, issuer, audience, and nonce (replay protection)
//  5. Extract username from preferred_username, name, or email claims
//  6. Determine user role (admin if in configured admin list, otherwise user)
//  7. Create encrypted session cookie with username and role
//  8. Clear state cookie and redirect to frontend
//
// Security measures:
//   - State cookie is always cleared (via defer) to prevent reuse
//   - ID token is cryptographically verified against the IdP's public keys
//   - Nonce in ID token must match the nonce from login to prevent replay attacks
//   - Session cookie is HttpOnly, Secure (in production), SameSite=Lax, encrypted
//   - Session expiration is configurable via auth.session_duration
func (h *OIDCAuthHandler) Callback(c *gin.Context) {
	defer h.clearCookie(c, stateCookieName)
	defer h.clearCookie(c, returnToCookieName)

	code := c.Query("code")
	stateParam := c.Query("state")

	if code == "" || stateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state parameter"})

		return
	}

	stateData, err := h.validateState(c, stateParam)
	if err != nil {
		return
	}

	oauth2Token, err := h.exchangeCode(c, code)
	if err != nil {
		return
	}

	idToken, err := h.verifyIDToken(c, oauth2Token, stateData.Nonce)
	if err != nil {
		return
	}

	username, err := h.extractUsername(c, idToken)
	if err != nil {
		return
	}

	role := h.determineRole(username)

	if err := h.createSession(c, username, role); err != nil {
		return
	}

	redirectPath := "/"

	if cookiePath, err := c.Cookie(returnToCookieName); err == nil && cookiePath != "" {
		if strings.HasPrefix(cookiePath, "/") && !strings.HasPrefix(cookiePath, "//") {
			redirectPath = cookiePath
		}
	}

	finalURL := strings.TrimRight(h.frontendURL, "/") + redirectPath

	slog.Info("OIDC authentication successful", "username", username, "role", role)

	c.Redirect(http.StatusFound, finalURL)
}

// Logout terminates the user's session by clearing the session cookie.
// The cookie is invalidated by setting MaxAge to -1, which instructs the browser
// to delete it immediately. This is a local logout only; the IdP session is not affected.
func (h *OIDCAuthHandler) Logout(c *gin.Context) {
	h.clearCookie(c, sessionCookieName)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func (h *OIDCAuthHandler) GetSessionManager() *session.Manager {
	return h.sessionMgr
}

func (h *OIDCAuthHandler) GetAdmins() []string {
	return h.admins
}

// validateState verifies the state parameter from the OIDC callback matches
// the state stored in the encrypted cookie. This prevents CSRF attacks by ensuring
// the callback request originated from a login flow initiated by this server.
//
// The state cookie is decrypted and validated, then the state value is compared
// to the query parameter. Returns the decoded state data (including nonce) on success.
func (h *OIDCAuthHandler) validateState(c *gin.Context, stateParam string) (*session.StateData, error) {
	stateCookie, err := c.Cookie(stateCookieName)
	if err != nil {
		slog.Error("Failed to get state cookie", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state cookie"})

		return nil, err
	}

	stateData, err := h.sessionMgr.DecodeState(stateCookie)
	if err != nil {
		slog.Error("Failed to decode state", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})

		return nil, err
	}

	if stateData.State != stateParam {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State mismatch"})

		return nil, fmt.Errorf("state mismatch")
	}

	return stateData, nil
}

// exchangeCode exchanges the authorization code for an OAuth2 token.
// This is a server-to-server call to the IdP's token endpoint.
func (h *OIDCAuthHandler) exchangeCode(c *gin.Context, code string) (*oauth2.Token, error) {
	oauth2Token, err := h.oauth2Config.Exchange(c.Request.Context(), code)
	if err != nil {
		slog.Error("Failed to exchange code for token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code"})

		return nil, err
	}

	return oauth2Token, nil
}

// verifyIDToken extracts and validates the ID token from the OAuth2 token response.
//
// Validation includes:
//   - Token signature verification using the IdP's public keys (JWKS)
//   - Issuer verification (must match configured issuer)
//   - Audience verification (must match configured client ID)
//   - Expiration check (token must not be expired)
//   - Nonce verification (must match the nonce from login to prevent replay attacks)
//
// Returns the verified ID token containing user claims on success.
func (h *OIDCAuthHandler) verifyIDToken(
	c *gin.Context,
	oauth2Token *oauth2.Token,
	nonce string,
) (*oidc.IDToken, error) {
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing id_token"})

		return nil, fmt.Errorf("missing id_token")
	}

	idToken, err := h.verifier.Verify(c.Request.Context(), rawIDToken)
	if err != nil {
		slog.Error("Failed to verify ID token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify id_token"})

		return nil, err
	}

	if idToken.Nonce != nonce {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nonce mismatch"})

		return nil, fmt.Errorf("nonce mismatch")
	}

	return idToken, nil
}

// extractUsername retrieves the username from the ID token claims.
// It tries claims in order: preferred_username, name, then email (local part).
// Returns an error if no suitable username claim is found.
func (h *OIDCAuthHandler) extractUsername(c *gin.Context, idToken *oidc.IDToken) (string, error) {
	var claims struct {
		PreferredUsername string `json:"preferred_username"`
		Name              string `json:"name"`
		Email             string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("Failed to parse claims", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})

		return "", err
	}

	username := claims.PreferredUsername
	if username == "" {
		username = claims.Name
	}

	if username == "" && claims.Email != "" {
		username = strings.Split(claims.Email, "@")[0]
	}

	if username == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not determine username"})

		return "", fmt.Errorf("could not determine username")
	}

	return username, nil
}

// determineRole assigns a role to the user based on the configured admin list.
// Returns "admin" if the username is in the admin list, otherwise "user".
func (h *OIDCAuthHandler) determineRole(username string) string {
	role := "user"

	for _, admin := range h.admins {
		if admin == username {
			role = "admin"

			break
		}
	}

	return role
}

// createSession generates an encrypted session cookie containing the user's
// username and role. The session data is encrypted using gorilla/securecookie
// with keys derived from the configured session secret.
//
// Security measures:
//   - Cookie is HttpOnly to prevent XSS attacks from accessing the session
//   - Cookie uses Secure flag in production (HTTPS only)
//   - Cookie uses SameSite=Lax to provide CSRF protection
//   - Session expiration is configurable via auth.session_duration
//   - Session data includes expiration timestamp for server-side validation
func (h *OIDCAuthHandler) createSession(c *gin.Context, username, role string) error {
	sessionDuration := h.sessionMgr.GetSessionDuration()

	sessionData := session.SessionData{
		Username:  username,
		Role:      role,
		ExpiresAt: time.Now().Add(sessionDuration),
	}

	encodedSession, err := h.sessionMgr.EncodeSession(sessionData)
	if err != nil {
		slog.Error("Failed to encode session", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})

		return err
	}

	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // Secure is set dynamically via h.secureCookie
		Name:     sessionCookieName,
		Value:    encodedSession,
		MaxAge:   int(sessionDuration.Seconds()),
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// clearCookie removes a cookie by setting its MaxAge to -1, which instructs
// the browser to delete it immediately. Used to clean up state and session cookies.
func (h *OIDCAuthHandler) clearCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // Secure is set dynamically via h.secureCookie
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
