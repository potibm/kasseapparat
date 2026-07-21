package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomString(t *testing.T) {
	str1, err := GenerateRandomString(RandomStringLength)
	require.NoError(t, err)
	assert.NotEmpty(t, str1)

	str2, err := GenerateRandomString(RandomStringLength)
	require.NoError(t, err)
	assert.NotEmpty(t, str2)

	assert.NotEqual(t, str1, str2, "Generated strings should be unique")
}

func TestSessionManager_EncodeDecode(t *testing.T) {
	mgr := NewManager("test-secret-key-that-is-at-least-32-characters-long", 24*time.Hour)

	data := SessionData{
		Username:  "testuser",
		Role:      "admin",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	encoded, err := mgr.EncodeSession(data)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := mgr.DecodeSession(encoded)
	require.NoError(t, err)
	assert.Equal(t, data.Username, decoded.Username)
	assert.Equal(t, data.Role, decoded.Role)
}

func TestSessionManager_RejectExpired(t *testing.T) {
	mgr := NewManager("test-secret-key-that-is-at-least-32-characters-long", 24*time.Hour)

	data := SessionData{
		Username:  "testuser",
		Role:      "user",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	encoded, err := mgr.EncodeSession(data)
	require.NoError(t, err)

	_, err = mgr.DecodeSession(encoded)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestSessionManager_RejectInvalid(t *testing.T) {
	mgr := NewManager("test-secret-key-that-is-at-least-32-characters-long", 24*time.Hour)

	_, err := mgr.DecodeSession("invalid-session-data")
	assert.Error(t, err)
}

func TestStateManager_EncodeDecode(t *testing.T) {
	mgr := NewManager("test-secret-key-that-is-at-least-32-characters-long", 24*time.Hour)

	data := StateData{
		State:     "test-state-value",
		Nonce:     "test-nonce-value",
		ExpiresAt: time.Now().Add(StateDuration),
	}

	encoded, err := mgr.EncodeState(data)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := mgr.DecodeState(encoded)
	require.NoError(t, err)
	assert.Equal(t, data.State, decoded.State)
	assert.Equal(t, data.Nonce, decoded.Nonce)
}

func TestStateManager_RejectExpired(t *testing.T) {
	mgr := NewManager("test-secret-key-that-is-at-least-32-characters-long", 24*time.Hour)

	data := StateData{
		State:     "test-state",
		Nonce:     "test-nonce",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	encoded, err := mgr.EncodeState(data)
	require.NoError(t, err)

	_, err = mgr.DecodeState(encoded)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestStateManager_RejectInvalid(t *testing.T) {
	mgr := NewManager("test-secret-key-that-is-at-least-32-characters-long", 24*time.Hour)

	_, err := mgr.DecodeState("invalid-state-data")
	assert.Error(t, err)
}

func TestDifferentManagersCannotDecodeEachOther(t *testing.T) {
	mgr1 := NewManager("secret-key-one-that-is-long-enough", 24*time.Hour)
	mgr2 := NewManager("secret-key-two-that-is-long-enough", 24*time.Hour)

	data := SessionData{
		Username:  "testuser",
		Role:      "user",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	encoded, err := mgr1.EncodeSession(data)
	require.NoError(t, err)

	_, err = mgr2.DecodeSession(encoded)
	assert.Error(t, err, "Different managers should not be able to decode each other's data")
}

func TestManager_GetSessionDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "24 hours",
			duration: 24 * time.Hour,
		},
		{
			name:     "1 hour",
			duration: 1 * time.Hour,
		},
		{
			name:     "30 minutes",
			duration: 30 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := NewManager("test-secret-key-that-is-at-least-32-characters-long", tt.duration)
			assert.Equal(t, tt.duration, mgr.GetSessionDuration())
		})
	}
}
