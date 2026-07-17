package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gorilla/securecookie"
)

const (
	SessionCookieName = "auth_session"
	StateCookieName   = "oidc_state"
	StateDuration     = 10 * time.Minute

	RandomStringLength = 32
)

type SessionData struct {
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

type StateData struct {
	State     string    `json:"state"`
	Nonce     string    `json:"nonce"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Manager struct {
	cookie          *securecookie.SecureCookie
	sessionDuration time.Duration
}

func NewManager(sessionSecret string, sessionDuration time.Duration) *Manager {
	hashKey := deriveKey(sessionSecret, "hash")
	blockKey := deriveKey(sessionSecret, "block")

	sc := securecookie.New(hashKey, blockKey)
	sc.MaxAge(int(sessionDuration.Seconds()))

	return &Manager{cookie: sc, sessionDuration: sessionDuration}
}

func (m *Manager) GetSessionDuration() time.Duration {
	return m.sessionDuration
}

func deriveKey(secret, purpose string) []byte {
	combined := secret + ":" + purpose
	hash := sha256.Sum256([]byte(combined))

	return hash[:]
}

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (m *Manager) EncodeSession(data SessionData) (string, error) {
	encoded, err := m.cookie.Encode(SessionCookieName, data)
	if err != nil {
		return "", fmt.Errorf("failed to encode session: %w", err)
	}

	return encoded, nil
}

func (m *Manager) DecodeSession(value string) (*SessionData, error) {
	var data SessionData
	if err := m.cookie.Decode(SessionCookieName, value, &data); err != nil {
		return nil, fmt.Errorf("failed to decode session: %w", err)
	}

	if time.Now().After(data.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	return &data, nil
}

func (m *Manager) EncodeState(data StateData) (string, error) {
	encoded, err := m.cookie.Encode(StateCookieName, data)
	if err != nil {
		return "", fmt.Errorf("failed to encode state: %w", err)
	}

	return encoded, nil
}

func (m *Manager) DecodeState(value string) (*StateData, error) {
	var data StateData
	if err := m.cookie.Decode(StateCookieName, value, &data); err != nil {
		return nil, fmt.Errorf("failed to decode state: %w", err)
	}

	if time.Now().After(data.ExpiresAt) {
		return nil, fmt.Errorf("state expired")
	}

	return &data, nil
}
