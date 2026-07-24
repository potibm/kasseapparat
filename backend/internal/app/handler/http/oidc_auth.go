// Package http provides HTTP handlers for the Kasseapparat application,
// including OIDC authentication for the BFF (Backend-for-Frontend) pattern.
//
// # OIDC Authentication Flow
//
// The OIDC handler implements a secure authentication flow using the Authorization Code Grant:
//
//  1. Login: Generates cryptographically secure state, nonce, and PKCE challenge,
//     stores them in an encrypted HttpOnly cookie, and redirects to the Identity Provider (IdP).
//
//  2. Callback: Validates the state parameter (CSRF protection), exchanges the authorization
//     code for tokens (using the PKCE verifier), verifies the ID token and nonce (replay attack protection),
//     extracts user information from claims, creates an encrypted session cookie, and redirects to frontend.
//
//  3. Logout: Clears the session cookie to terminate the user's session.
//
// # Security Considerations
//
//   - CSRF Protection: The state parameter is generated using crypto/rand and validated on callback
//   - Replay Attack Prevention: The nonce is included in the ID token and verified
//   - PKCE (Proof Key for Code Exchange): Prevents authorization code interception attacks
//   - Secure Cookies: Session and state cookies are HttpOnly (XSS protection), use SameSite=Lax,
//     and Secure flag in production (HTTPS only)
//   - Encrypted Sessions: Session data is encrypted using gorilla/securecookie with derived keys
//   - Configurable Duration: Session lifetime is configurable via auth.session_duration
//   - Admin Detection: Users with the configured admin group claim receive "admin" role, others get "user"
package http

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
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

	pkceVerifierBytes = 32
)

type OIDCOptions struct {
	Issuer          string
	ClientID        string
	ClientSecret    string
	CallbackURL     string
	FrontendURL     string
	SessionSecret   string
	SessionDuration time.Duration
	AdminGroup      string
	IsProduction    bool
}

type OIDCAuthHandler struct {
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	sessionMgr   *session.Manager
	adminGroup   string
	frontendURL  string
	callbackURL  string
	secureCookie bool
	log          slog.Logger
}

// NewOIDCAuthHandler initializes the OIDC authentication handler with the provided options.
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
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: opts.ClientID})
	sessionMgr := session.NewManager(opts.SessionSecret, opts.SessionDuration)

	return &OIDCAuthHandler{
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		sessionMgr:   sessionMgr,
		adminGroup:   opts.AdminGroup,
		frontendURL:  opts.FrontendURL,
		callbackURL:  opts.CallbackURL,
		secureCookie: opts.IsProduction,
	}, nil
}

// Login initiates the OIDC authentication flow by generating a secure state,
// nonce, and PKCE challenge, storing them in an encrypted cookie, and redirecting to the IdP.
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

	// PKCE Code Verifier & Challenge
	verifierBytes := make([]byte, pkceVerifierBytes)
	if _, err := rand.Read(verifierBytes); err != nil {
		slog.Error("Failed to generate code_verifier", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PKCE verifier"})

		return
	}

	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes)

	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	stateData := session.StateData{
		State:        state,
		Nonce:        nonce,
		CodeVerifier: codeVerifier,
		ExpiresAt:    time.Now().Add(session.StateDuration),
	}

	encodedState, err := h.sessionMgr.EncodeState(stateData)
	if err != nil {
		slog.Error("Failed to encode state", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode state"})

		return
	}

	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // HttpOnly cookie for CSRF protection
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

	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // HttpOnly cookie for return URL
		Name:     returnToCookieName,
		Value:    returnTo,
		MaxAge:   int(session.StateDuration.Seconds()),
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// PKCE parameters for the redirect
	authURL := h.oauth2Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	slog.Debug("Redirecting to OIDC provider", "url", authURL)

	c.Redirect(http.StatusFound, authURL)
}

// Callback handles the OIDC callback from the identity provider.
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

	// Exchange code with PKCE verifier
	oauth2Token, err := h.exchangeCode(c, code, stateData.CodeVerifier)
	if err != nil {
		return
	}

	idToken, err := h.verifyIDToken(c, oauth2Token, stateData.Nonce)
	if err != nil {
		return
	}

	username, groups, err := h.extractUserClaims(c, idToken)
	if err != nil {
		return
	}

	role := h.determineRole(groups)

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
func (h *OIDCAuthHandler) Logout(c *gin.Context) {
	h.clearCookie(c, sessionCookieName)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func (h *OIDCAuthHandler) GetSessionManager() *session.Manager {
	return h.sessionMgr
}

func (h *OIDCAuthHandler) GetAdminGroup() string {
	return h.adminGroup
}

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

// exchangeCode exchanges the authorization code for an OAuth2 token using PKCE.
func (h *OIDCAuthHandler) exchangeCode(c *gin.Context, code, codeVerifier string) (*oauth2.Token, error) {
	oauth2Token, err := h.oauth2Config.Exchange(
		c.Request.Context(),
		code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		slog.Error("Failed to exchange code for token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code"})

		return nil, err
	}

	return oauth2Token, nil
}

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

func (h *OIDCAuthHandler) extractUserClaims(
	c *gin.Context,
	idToken *oidc.IDToken,
) (username string, groups []string, err error) {
	var claims struct {
		PreferredUsername string   `json:"preferred_username"`
		Name              string   `json:"name"`
		Email             string   `json:"email"`
		Groups            []string `json:"groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("Failed to parse claims", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})

		return "", nil, err
	}

	username = claims.PreferredUsername
	if username == "" {
		username = claims.Name
	}

	if username == "" && claims.Email != "" {
		username = strings.Split(claims.Email, "@")[0]
	}

	if username == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not determine username"})

		return "", nil, fmt.Errorf("could not determine username")
	}

	return username, claims.Groups, nil
}

func (h *OIDCAuthHandler) determineRole(groups []string) string {
	if h.adminGroup != "" && slices.Contains(groups, h.adminGroup) {
		return "admin"
	}

	return "user"
}

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

	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // HttpOnly session cookie
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

func (h *OIDCAuthHandler) clearCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // Cookie clearing is safe
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   h.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
