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
	stateCookieName   = "oidc_state"
	sessionCookieName = "auth_session"
)

type OIDCAuthHandler struct {
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	sessionMgr   *session.Manager
	admins       []string
	frontendURL  string
	callbackURL  string
}

func NewOIDCAuthHandler(
	ctx context.Context,
	issuer, clientID, clientSecret, callbackURL, frontendURL, sessionSecret string,
	admins []string,
) (*OIDCAuthHandler, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	sessionMgr := session.NewManager(sessionSecret)

	return &OIDCAuthHandler{
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		sessionMgr:   sessionMgr,
		admins:       admins,
		frontendURL:  frontendURL,
		callbackURL:  callbackURL,
	}, nil
}

func (h *OIDCAuthHandler) Login(c *gin.Context) {
	state, err := session.GenerateRandomString(session.RandomStringLength)
	if err != nil {
		slog.Error("Failed to generate state", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})

		return
	}

	nonce, err := session.GenerateRandomString(session.RandomStringLength)
	if err != nil {
		slog.Error("Failed to generate nonce", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate nonce"})

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode state"})

		return
	}

	c.SetCookie(stateCookieName, encodedState, int(session.StateDuration.Seconds()), "/", "", false, true)

	authURL := h.oauth2Config.AuthCodeURL(state) + "&nonce=" + nonce

	slog.Debug("Redirecting to OIDC provider", "url", authURL)

	c.Redirect(http.StatusFound, authURL)
}

func (h *OIDCAuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	stateParam := c.Query("state")

	if code == "" || stateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code or state parameter"})

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

	c.SetCookie(stateCookieName, "", -1, "/", "", false, true)

	slog.Info("OIDC authentication successful", "username", username, "role", role)

	c.Redirect(http.StatusFound, h.frontendURL)
}

func (h *OIDCAuthHandler) Logout(c *gin.Context) {
	c.SetCookie(sessionCookieName, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *OIDCAuthHandler) GetSessionManager() *session.Manager {
	return h.sessionMgr
}

func (h *OIDCAuthHandler) GetAdmins() []string {
	return h.admins
}

func (h *OIDCAuthHandler) validateState(c *gin.Context, stateParam string) (*session.StateData, error) {
	stateCookie, err := c.Cookie(stateCookieName)
	if err != nil {
		slog.Error("Failed to get state cookie", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing state cookie"})

		return nil, err
	}

	stateData, err := h.sessionMgr.DecodeState(stateCookie)
	if err != nil {
		slog.Error("Failed to decode state", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})

		return nil, err
	}

	if stateData.State != stateParam {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state mismatch"})

		return nil, fmt.Errorf("state mismatch")
	}

	return stateData, nil
}

func (h *OIDCAuthHandler) exchangeCode(c *gin.Context, code string) (*oauth2.Token, error) {
	oauth2Token, err := h.oauth2Config.Exchange(c.Request.Context(), code)
	if err != nil {
		slog.Error("Failed to exchange code for token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code"})

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "missing id_token"})

		return nil, fmt.Errorf("missing id_token")
	}

	idToken, err := h.verifier.Verify(c.Request.Context(), rawIDToken)
	if err != nil {
		slog.Error("Failed to verify ID token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify id_token"})

		return nil, err
	}

	if idToken.Nonce != nonce {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nonce mismatch"})

		return nil, fmt.Errorf("nonce mismatch")
	}

	return idToken, nil
}

func (h *OIDCAuthHandler) extractUsername(c *gin.Context, idToken *oidc.IDToken) (string, error) {
	var claims struct {
		PreferredUsername string `json:"preferred_username"`
		Name              string `json:"name"`
		Email             string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("Failed to parse claims", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse claims"})

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not determine username"})

		return "", fmt.Errorf("could not determine username")
	}

	return username, nil
}

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

func (h *OIDCAuthHandler) createSession(c *gin.Context, username, role string) error {
	sessionData := session.SessionData{
		Username:  username,
		Role:      role,
		ExpiresAt: time.Now().Add(session.SessionDuration),
	}

	encodedSession, err := h.sessionMgr.EncodeSession(sessionData)
	if err != nil {
		slog.Error("Failed to encode session", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})

		return err
	}

	c.SetCookie(sessionCookieName, encodedSession, int(session.SessionDuration.Seconds()), "/", "", false, true)

	return nil
}
