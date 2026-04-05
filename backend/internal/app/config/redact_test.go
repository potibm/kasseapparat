package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactConfigForDisplay(t *testing.T) {
	// Initialize a config with sensitive data
	cfg := Config{}
	cfg.Jwt.Secret = "super-secret-key"
	cfg.Sumup.ApiKey = "sup_sk_12345"
	cfg.Sentry.DSN = "https://public@sentry.io/1"
	cfg.App.RedisURL = "redis://:p@ssword@localhost:6379/0"
	cfg.Mailer.DSN = "smtp://user:secret-mail-pass@smtp.example.com:587"

	// Execute redaction
	redactedCfg := cfg.RedactConfigForDisplay()
	urlEncodedRedacted := "%2A%2A%2AREDACTED%2A%2A%2"

	// Verify standard fields are redacted
	assert.Equal(t, redacted, redactedCfg.Jwt.Secret)
	assert.Equal(t, redacted, redactedCfg.Sumup.ApiKey)
	assert.Equal(t, redacted, redactedCfg.Sentry.DSN)

	// Verify URL passwords are redacted but structure remains
	assert.Contains(t, string(redactedCfg.App.RedisURL), urlEncodedRedacted)
	assert.Contains(t, redactedCfg.Mailer.DSN, urlEncodedRedacted)
	assert.Contains(t, redactedCfg.Mailer.DSN, "user") // Username should still be visible
	assert.NotContains(t, redactedCfg.Mailer.DSN, "secret-mail-pass")

	// Ensure the original config was not modified (Immutability check)
	assert.Equal(t, "super-secret-key", cfg.Jwt.Secret)
}

func TestRedactUrlPassword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with user and password",
			input:    "postgres://admin:password123@localhost:5432/db",
			expected: "postgres://admin:%2A%2A%2AREDACTED%2A%2A%2A@localhost:5432/db",
		},
		{
			name:     "URL with password only",
			input:    "redis://:onlypass@localhost:6379",
			expected: "redis://:%2A%2A%2AREDACTED%2A%2A%2A@localhost:6379",
		},
		{
			name:     "URL without credentials",
			input:    "https://example.com/api",
			expected: "https://example.com/api",
		},
		{
			name:     "Invalid URL",
			input:    "://invalid-url",
			expected: "://invalid-url", // Should return raw input on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := redactUrlPassword(tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}
