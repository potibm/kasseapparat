package config

import (
	"log"
	"os"
	"strings"
)

type SentryConfig struct {
	DSN                     string
	TraceSampleRate         float64
	ReplaySessionSampleRate float64
	ReplayErrorSampleRate   float64
	Environment             string
}

type JwtConfig struct {
	Secret string
	Realm  string
}

type MailConfig struct {
	DSN               string
	FromEmail         string
	MailSubjectPrefix string
}

type AppConfig struct {
	Version string
	GinMode string
}

type FormatConfig struct {
	CurrencyLocale    string
	CurrencyCode      string
	DateLocale        string
	DateOptions       string
	FractionDigitsMin int
	FractionDigitsMax int
}

type Config struct {
	AppConfig          AppConfig
	FormatConfig       FormatConfig
	VATRates           string
	EnvironmentMessage string
	PaymentMethods     PaymentMethods
	SentryConfig       SentryConfig
	JwtConfig          JwtConfig
	CorsAllowOrigins   []string
	FrontendURL        string
	MailConfig         MailConfig
}

func Load() Config {
	return Config{
		AppConfig:          loadAppConfig(),
		FormatConfig:       loadFormatConfig(),
		VATRates:           getEnvWithJSONValidation("VAT_RATES", "[{\"rate\":25,\"name\":\"Standard\"},{\"rate\":0,\"name\":\"Zero rate\"}]"),
		EnvironmentMessage: getEnv("ENV_MESSAGE", ""),
		PaymentMethods:     loadPaymentMethods(),
		SentryConfig:       loadSentryConfig(),
		JwtConfig:          loadJwtConfig(),
		CorsAllowOrigins:   loadCorsAllowOrigins(),
		FrontendURL:        getEnv("FRONTEND_URL", ""),
		MailConfig:         loadMailConfig(),
	}
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Version: readVersionFromFile(),
		GinMode: getEnv("GIN_MODE", "release"),
	}
}

func readVersionFromFile() string {
	versionFilePath := "./VERSION"

	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		log.Printf("Error reading the version file: %v", err)

		return "0.0.0"
	}

	return strings.TrimSpace(string(content))
}

func loadCorsAllowOrigins() []string {
	origins := getEnv("CORS_ALLOW_ORIGINS", "")
	if origins == "" {
		log.Fatalf("CORS_ALLOW_ORIGINS is not set in env")
	}

	return strings.Split(origins, ",")
}

func loadFormatConfig() FormatConfig {
	return FormatConfig{
		CurrencyLocale:    getEnv("CURRENCY_LOCALE", "dk-DK"),
		CurrencyCode:      getEnv("CURRENCY_CODE", "DKK"),
		DateLocale:        getEnv("DATE_LOCALE", "dk-DK"),
		DateOptions:       getEnvWithJSONValidation("DATE_OPTIONS", "{\"weekday\":\"long\",\"hour\":\"2-digit\",\"minute\":\"2-digit\"}"),
		FractionDigitsMin: getEnvAsInt("FRACTION_DIGITS_MIN", 0),
		FractionDigitsMax: getEnvAsInt("FRACTION_DIGITS_MAX", 2),
	}
}

func loadSentryConfig() SentryConfig {
	return SentryConfig{
		DSN:                     getEnv("SENTRY_DSN", ""),
		TraceSampleRate:         getEnvAsFloat("SENTRY_TRACE_SAMPLE_RATE", 0.1),
		ReplaySessionSampleRate: getEnvAsFloat("SENTRY_REPLAY_SESSION_SAMPLE_RATE", 0.1),
		ReplayErrorSampleRate:   getEnvAsFloat("SENTRY_REPLAY_ERROR_SAMPLE_RATE", 0.1),
	}
}

func loadJwtConfig() JwtConfig {
	return JwtConfig{
		Realm:  getEnv("JWT_REALM", "kasseapparat"),
		Secret: getEnv("JWT_SECRET", ""),
	}
}

func loadMailConfig() MailConfig {
	return MailConfig{
		DSN:               getEnv("MAIL_DSN", ""),
		FromEmail:         getEnv("MAIL_FROM", ""),
		MailSubjectPrefix: getEnv("MAIL_SUBJECT_PREFIX", "[Kasseapparat]"),
	}
}
