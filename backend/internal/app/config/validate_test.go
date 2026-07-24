package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultTestConfig = Config{
	App: AppConfig{
		DbFilename:  "kasseapparat",
		GinMode:     "release",
		Environment: "production",
		LogLevel:    "info",
		LogFormat:   "json",
		FrontendURL: "http://localhost:3000",
		Port:        8080,
	},
	Format: FormatConfig{
		Currency: CurrencyFormatConfig{Locale: "de-DE", Code: "EUR"},
		Date:     DateFormatConfig{Locale: "en-US"},
	},
	Mailer: MailerConfig{
		DSN:               "smtp://user:pass@localhost:587",
		FromEmail:         "noreply@example.com",
		MailSubjectPrefix: "[Kass]",
		FrontendURL:       "http://localhost:3000",
	},
	Sentry: SentryConfig{DSN: ""},
	Auth: AuthConfig{
		Mode:        "proxy",
		ProxyHeader: "X-Remote-User",
		ProxyAdmins: []string{},
	},
}

func TestConfigValidate(t *testing.T) {
	cfg := defaultTestConfig
	assert.NoError(t, cfg.Validate())
}

func TestCurrencyFormatConfigValidate(t *testing.T) {
	cfg := CurrencyFormatConfig{Locale: "de-DE", Code: "EUR"}
	assert.NoError(t, cfg.Validate())

	cfg = CurrencyFormatConfig{Locale: "xx", Code: "EUR"}
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currency.locale 'xx' is not a valid locale")

	cfg = CurrencyFormatConfig{Locale: "de-DE", Code: "XX"}
	err = cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currency.code 'XX' is not a valid ISO 4217 code")
}

func TestDateFormatConfigValidate(t *testing.T) {
	cfg := DateFormatConfig{Locale: "en-US"}
	assert.NoError(t, cfg.Validate())

	cfg = DateFormatConfig{Locale: "invalid-locale"}
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "date.locale 'invalid-locale' is not a valid locale")
}

func TestAppConfigValidate(t *testing.T) {
	cfg := AppConfig{DbFilename: "kasseapparat"}
	assert.NoError(t, cfg.Validate())

	cfg = AppConfig{DbFilename: ""}
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db_filename '' contains invalid characters")

	cfg = AppConfig{DbFilename: "../invalid"}
	err = cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db_filename '../invalid' contains invalid characters")
}

func TestAuthConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      AuthConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid proxy config",
			config: AuthConfig{
				Mode:        "proxy",
				ProxyHeader: "X-Remote-User",
			},
			expectError: false,
		},
		{
			name: "valid oidc config",
			config: AuthConfig{
				Mode:             "oidc",
				OidcIssuer:       "https://auth.example.com",
				OidcClientID:     "client-id",
				OidcClientSecret: "client-secret",
				SessionSecret:    "this-is-a-very-long-secret-that-is-at-least-32-chars",
			},
			expectError: false,
		},
		{
			name: "proxy mode without header",
			config: AuthConfig{
				Mode:        "proxy",
				ProxyHeader: "",
			},
			expectError: true,
			errorMsg:    "auth.proxy_header is required when mode is 'proxy'",
		},
		{
			name: "oidc mode without issuer",
			config: AuthConfig{
				Mode:             "oidc",
				OidcIssuer:       "",
				OidcClientID:     "client-id",
				OidcClientSecret: "client-secret",
				SessionSecret:    "this-is-a-very-long-secret-that-is-at-least-32-chars",
			},
			expectError: true,
			errorMsg:    "auth.oidc_issuer is required when mode is 'oidc'",
		},
		{
			name: "oidc mode without client id",
			config: AuthConfig{
				Mode:             "oidc",
				OidcIssuer:       "https://auth.example.com",
				OidcClientID:     "",
				OidcClientSecret: "client-secret",
				SessionSecret:    "this-is-a-very-long-secret-that-is-at-least-32-chars",
			},
			expectError: true,
			errorMsg:    "auth.oidc_client_id is required when mode is 'oidc'",
		},
		{
			name: "oidc mode without client secret",
			config: AuthConfig{
				Mode:             "oidc",
				OidcIssuer:       "https://auth.example.com",
				OidcClientID:     "client-id",
				OidcClientSecret: "",
				SessionSecret:    "this-is-a-very-long-secret-that-is-at-least-32-chars",
			},
			expectError: true,
			errorMsg:    "auth.oidc_client_secret is required when mode is 'oidc'",
		},
		{
			name: "oidc mode with short session secret",
			config: AuthConfig{
				Mode:             "oidc",
				OidcIssuer:       "https://auth.example.com",
				OidcClientID:     "client-id",
				OidcClientSecret: "client-secret",
				SessionSecret:    "too-short",
			},
			expectError: true,
			errorMsg:    "auth.session_secret must be at least 32 characters when mode is 'oidc'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
