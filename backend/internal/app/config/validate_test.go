package config

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultTestConfig = Config{
	App: AppConfig{
		DbFilename:  "kasseapparat",
		RedisURL:    "",
		GinMode:     "release",
		Environment: "production",
		LogLevel:    "info",
		LogFormat:   "json",
		FrontendURL: "http://localhost:3000",
	},
	Format: FormatConfig{
		Currency: CurrencyFormatConfig{Locale: "de-DE", Code: "EUR"},
		Date:     DateFormatConfig{Locale: "en-US"},
	},
	Jwt: JwtConfig{Secret: "asecretforsec", Realm: "kasseapparat"},
	Mailer: MailerConfig{
		DSN:               "smtp://user:pass@localhost:587",
		FromEmail:         "noreply@example.com",
		MailSubjectPrefix: "[Kass]",
		FrontendURL:       "http://localhost:3000",
	},
	Sentry: SentryConfig{DSN: ""},
}

func TestConfigValidate(t *testing.T) {
	var buf bytes.Buffer

	h := slog.NewJSONHandler(&buf, nil)
	logger := slog.New(h)
	oldLogger := slog.Default()

	slog.SetDefault(logger)
	t.Cleanup(func() { slog.SetDefault(oldLogger) })

	cfg := defaultTestConfig
	assert.NoError(t, cfg.Validate())

	assert.NotContains(t, buf.String(), "WARN", "Expected no warnings when using a non-default JWT secret")
}

func TestConfigValidateWithDefaultJwtSecretReturningErrorInProduction(t *testing.T) {
	cfg := defaultTestConfig
	cfg.App.Environment = "production"
	cfg.Jwt.Secret = DefaultJwtSecret

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET is set to the default value, which is not allowed in production")
}

func TestConfigValidateWithDefaultJwtSecretShowingWarning(t *testing.T) {
	var buf bytes.Buffer

	h := slog.NewJSONHandler(&buf, nil)
	logger := slog.New(h)
	oldLogger := slog.Default()

	slog.SetDefault(logger)
	t.Cleanup(func() { slog.SetDefault(oldLogger) })

	cfg := defaultTestConfig
	cfg.App.Environment = "development"
	cfg.Jwt.Secret = DefaultJwtSecret

	err := cfg.Validate()
	assert.NoError(t, err)

	assert.Contains(
		t,
		buf.String(),
		"\"level\":\"WARN\"",
		"The expected warning about using the default JWT secret was not logged",
	)
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

func TestRedisUrlValidate(t *testing.T) {
	validURL := RedisURL("redis://user:password@localhost:6379/0")
	assert.NoError(t, validURL.Validate())

	invalidURL := RedisURL("not-a-valid-url")
	err := invalidURL.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis_url 'not-a-valid-url' is not a valid URL")

	invalidScheme := RedisURL("http://localhost:6379")
	err = invalidScheme.Validate()
	assert.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"redis_url 'http://localhost:6379' has invalid scheme 'http' (expected 'redis' or 'rediss')",
	)

	missingHost := RedisURL("redis:///0")
	err = missingHost.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis_url 'redis:///0' has missing host")
}

func TestAppConfigValidate(t *testing.T) {
	cfg := AppConfig{DbFilename: "kasseapparat", RedisURL: ""}
	assert.NoError(t, cfg.Validate())

	cfg = AppConfig{DbFilename: "", RedisURL: ""}
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db_filename '' contains invalid characters")

	cfg = AppConfig{DbFilename: "../invalid", RedisURL: ""}
	err = cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db_filename '../invalid' contains invalid characters")
}
