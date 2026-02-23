package config

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/joho/godotenv"
)

const (
	DefaultVatRates                = "[{\"rate\":25,\"name\":\"Standard\"},{\"rate\":0,\"name\":\"Zero rate\"}]"
	DefaultDateOptions             = "{\"weekday\":\"long\",\"hour\":\"2-digit\",\"minute\":\"2-digit\"}"
	defaultTraceSampleRate         = 0.1
	defaultReplaySessionSampleRate = 0.1
	defaultReplayErrorSampleRate   = 0.1
	defaultMinorUnit               = 2
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
	Secret       string
	Realm        string
	SecureCookie bool
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

type CorsAllowOriginsConfig []string

type Config struct {
	AppConfig          AppConfig
	FormatConfig       FormatConfig
	VATRates           string
	EnvironmentMessage string
	PaymentMethods     PaymentMethods
	SentryConfig       SentryConfig
	JwtConfig          JwtConfig
	CorsAllowOrigins   CorsAllowOriginsConfig
	FrontendURL        string
	MailerConfig       MailerConfig
	SumupConfig        SumupConfig
}

func (cfg Config) OutputVersion() {
	slog.Info("Kasseapparat", slog.String("version", cfg.AppConfig.Version))
}

func Load(logger *slog.Logger) (Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Error loading .env file, using environment variables")
	}

	config, err := loadConfig(logger)
	if err != nil {
		return Config{}, err
	}

	return *config, nil
}

func loadConfig(logger *slog.Logger) (*Config, error) {
	corsAllowOriginsConfig, err := loadCorsAllowOrigins()
	if err != nil {
		return nil, err
	}

	return &Config{
		AppConfig:          loadAppConfig(),
		FormatConfig:       loadFormatConfig(logger),
		VATRates:           getEnvWithJSONValidation(logger, "VAT_RATES", DefaultVatRates),
		EnvironmentMessage: getEnv("ENV_MESSAGE", ""),
		PaymentMethods:     loadPaymentMethods(),
		SentryConfig:       loadSentryConfig(),
		JwtConfig:          loadJwtConfig(),
		CorsAllowOrigins:   corsAllowOriginsConfig,
		MailerConfig:       loadMailerConfig(),
		SumupConfig:        loadSumupConfig(),
	}, nil
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Version: "0.0.0",
		GinMode: getEnv("GIN_MODE", "release"),
	}
}

func (cfg *Config) SetVersion(version string) {
	cfg.AppConfig.Version = version
	cfg.SentryConfig.Version = version
}

func loadCorsAllowOrigins() ([]string, error) {
	origins := getEnv("CORS_ALLOW_ORIGINS", "")
	if origins == "" {
		return nil, errors.New("CORS_ALLOW_ORIGINS is not set in env")
	}

	return strings.Split(origins, ","), nil
}

func loadFormatConfig(logger *slog.Logger) FormatConfig {
	return FormatConfig{
		CurrencyLocale:    getEnv("CURRENCY_LOCALE", "dk-DK"),
		CurrencyCode:      getCurrencyCode(),
		DateLocale:        getEnv("DATE_LOCALE", "dk-DK"),
		DateOptions:       getEnvWithJSONValidation(logger, "DATE_OPTIONS", DefaultDateOptions),
		FractionDigitsMin: getEnvAsInt("FRACTION_DIGITS_MIN", 0),
		FractionDigitsMax: getCurrencyMinorUnit(),
	}
}

func loadSentryConfig() SentryConfig {
	return SentryConfig{
		DSN:                     getEnv("SENTRY_DSN", ""),
		TraceSampleRate:         getEnvAsFloat("SENTRY_TRACE_SAMPLE_RATE", defaultTraceSampleRate),
		ReplaySessionSampleRate: getEnvAsFloat("SENTRY_REPLAY_SESSION_SAMPLE_RATE", defaultReplaySessionSampleRate),
		ReplayErrorSampleRate:   getEnvAsFloat("SENTRY_REPLAY_ERROR_SAMPLE_RATE", defaultReplayErrorSampleRate),
		Version:                 "0.0.0",
	}
}

func loadJwtConfig() JwtConfig {
	return JwtConfig{
		Realm:        getEnv("JWT_REALM", "kasseapparat"),
		Secret:       getEnv("JWT_SECRET", ""),
		SecureCookie: getEnvAsBool("JWT_SECURE_COOKIE", true),
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
	return getEnvAsInt("FRACTION_DIGITS_MAX", defaultMinorUnit)
}
