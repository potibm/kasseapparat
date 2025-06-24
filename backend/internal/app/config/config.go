package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type SentryConfig struct {
	DSN                     string
	TraceSampleRate         float64
	ReplaySessionSampleRate float64
	ReplayErrorSampleRate   float64
	Environment             string
	Version                 string
}

type JwtConfig struct {
	Secret string
	Realm  string
}

type MailerConfig struct {
	DSN               string
	FromEmail         string
	MailSubjectPrefix string
	FrontendURL       string
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
	MailerConfig       MailerConfig
	SumupConfig        SumupConfig
}

func (cfg Config) OutputVersion() {
	log.Printf("Kasseapparat %s\n", cfg.AppConfig.Version)
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	return Config{
		AppConfig:          loadAppConfig(),
		FormatConfig:       loadFormatConfig(),
		VATRates:           getEnvWithJSONValidation("VAT_RATES", "[{\"rate\":25,\"name\":\"Standard\"},{\"rate\":0,\"name\":\"Zero rate\"}]"),
		EnvironmentMessage: getEnv("ENV_MESSAGE", ""),
		PaymentMethods:     loadPaymentMethods(),
		SentryConfig:       loadSentryConfig(),
		JwtConfig:          loadJwtConfig(),
		CorsAllowOrigins:   loadCorsAllowOrigins(),
		MailerConfig:       loadMailerConfig(),
		SumupConfig:        loadSumupConfig(),
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
		CurrencyCode:      getCurrencyCode(),
		DateLocale:        getEnv("DATE_LOCALE", "dk-DK"),
		DateOptions:       getEnvWithJSONValidation("DATE_OPTIONS", "{\"weekday\":\"long\",\"hour\":\"2-digit\",\"minute\":\"2-digit\"}"),
		FractionDigitsMin: getEnvAsInt("FRACTION_DIGITS_MIN", 0),
		FractionDigitsMax: getCurrencyMinorUnit(),
	}
}

func loadSentryConfig() SentryConfig {
	return SentryConfig{
		DSN:                     getEnv("SENTRY_DSN", ""),
		TraceSampleRate:         getEnvAsFloat("SENTRY_TRACE_SAMPLE_RATE", 0.1),
		ReplaySessionSampleRate: getEnvAsFloat("SENTRY_REPLAY_SESSION_SAMPLE_RATE", 0.1),
		ReplayErrorSampleRate:   getEnvAsFloat("SENTRY_REPLAY_ERROR_SAMPLE_RATE", 0.1),
		Version:                 readVersionFromFile(),
	}
}

func loadJwtConfig() JwtConfig {
	return JwtConfig{
		Realm:  getEnv("JWT_REALM", "kasseapparat"),
		Secret: getEnv("JWT_SECRET", ""),
	}
}

func loadMailerConfig() MailerConfig {
	return MailerConfig{
		DSN:               getEnv("MAIL_DSN", "smtp://user:password@localhost:1025"),
		FromEmail:         getEnv("MAIL_FROM", ""),
		MailSubjectPrefix: getEnv("MAIL_SUBJECT_PREFIX", "[Kasseapparat]"),
		FrontendURL:       getEnv("FRONTEND_URL", ""),
	}
}

func getCurrencyCode() string {
	return getEnv("CURRENCY_CODE", "DKK")
}

func getCurrencyMinorUnit() int {
	return getEnvAsInt("FRACTION_DIGITS_MAX", 2)
}
