package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	cfg := Config{
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

func TestRedisUrlValidate(t *testing.T) {
	validUrl := RedisUrl("redis://user:password@localhost:6379/0")
	assert.NoError(t, validUrl.Validate())

	invalidUrl := RedisUrl("not-a-valid-url")
	err := invalidUrl.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis_url 'not-a-valid-url' is not a valid URL")

	invalidScheme := RedisUrl("http://localhost:6379")
	err = invalidScheme.Validate()
	assert.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		"redis_url 'http://localhost:6379' has invalid scheme 'http' (expected 'redis' or 'rediss')",
	)

	missingHost := RedisUrl("redis:///0")
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
