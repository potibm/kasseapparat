package config

import (
	"net/url"

	"github.com/potibm/kasseapparat/internal/app/models"
)

type SentryConfig struct {
	DSN                     string  `mapstructure:"dsn"                        validate:"omitempty,url"`
	TraceSampleRate         float64 `mapstructure:"trace_sample_rate"          validate:"omitempty,gte=0,lte=1"`
	ReplaySessionSampleRate float64 `mapstructure:"replay_session_sample_rate" validate:"omitempty,gte=0,lte=1"`
	ReplayErrorSampleRate   float64 `mapstructure:"replay_error_sample_rate"   validate:"omitempty,gte=0,lte=1"`
	Environment             string  `mapstructure:"environment"`
	Version                 string  `mapstructure:"version"`
}

type JwtConfig struct {
	Secret       string `mapstructure:"secret"        validate:"required,min=8"`
	Realm        string `mapstructure:"realm"         validate:"required"`
	SecureCookie bool   `mapstructure:"secure_cookie"`
}

type MailerConfig struct {
	DSN               string `mapstructure:"dsn"            validate:"required,url"`
	FromEmail         string `mapstructure:"from"           validate:"email"`
	MailSubjectPrefix string `mapstructure:"subject_prefix" validate:"required"`
	FrontendURL       string `mapstructure:"frontend_url"   validate:"required,http_url"`
}

type RedisUrl string

type RedisConfig url.URL

type AppConfig struct {
	Version string `mapstructure:"version"`

	GinMode     string `mapstructure:"gin_mode" validate:"required,oneof=debug release test"`
	Environment string `mapstructure:"env"      validate:"required,oneof=development staging production test"`

	LogLevel  string `mapstructure:"log_level"  validate:"required,oneof=debug info warn error"`
	LogFormat string `mapstructure:"log_format" validate:"required,oneof=json text"`

	DbFilename         string                 `mapstructure:"db_filename"         validate:"required"`
	RedisURL           RedisUrl               `mapstructure:"redis_url"           validate:"omitempty,url"`
	FrontendURL        string                 `mapstructure:"frontend_url"        validate:"required,http_url"`
	CorsAllowOrigins   CorsAllowOriginsConfig `mapstructure:"cors_allow_origins"  validate:"dive,required"`
	EnvironmentMessage string                 `mapstructure:"environment_message"`
}

type FormatConfig struct {
	Currency CurrencyFormatConfig `mapstructure:"currency"`
	Date     DateFormatConfig     `mapstructure:"date"`
}

type CurrencyFormatConfig struct {
	Locale            string `mapstructure:"locale"              validate:"required"`
	Code              string `mapstructure:"code"                validate:"required"`
	FractionDigitsMin int    `mapstructure:"fraction_digits_min" validate:"gte=0"`
	FractionDigitsMax int    `mapstructure:"fraction_digits_max" validate:"gte=0"`
}

type DateFormatOptionsConfig map[string]any

type DateFormatConfig struct {
	Locale  string                  `mapstructure:"locale"  validate:"required"`
	Options DateFormatOptionsConfig `mapstructure:"options"`
}

type CorsAllowOriginsConfig []string

type VatRateConfig struct {
	Rate float64 `mapstructure:"rate"`
	Name string  `mapstructure:"name"`
}

type VatRatesConfig []VatRateConfig

type PaymentMethods []PaymentMethodConfig

type PaymentMethodConfig struct {
	Code models.PaymentMethod
	Name string
}

type SumupConfig struct {
	ApiKey            string `mapstructure:"api_key"`
	MerchantCode      string `mapstructure:"merchant_code"`
	CurrencyCode      string `mapstructure:"currency_code"`
	CurrencyMinorUnit int    `mapstructure:"currency_minor_unit"`
	AffiliateKey      string `mapstructure:"affiliate_key"`
	ApplicationId     string `mapstructure:"application_id"`
	PublicUrl         string `mapstructure:"public_url"          validate:"omitempty,https_url"`
}

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Format FormatConfig `mapstructure:"format"`
	Sentry SentryConfig `mapstructure:"sentry"`
	Jwt    JwtConfig    `mapstructure:"jwt"`
	Mailer MailerConfig `mapstructure:"mailer"`
	Sumup  SumupConfig  `mapstructure:"sumup"`

	VATRates       VatRatesConfig `mapstructure:"vat_rates"`
	PaymentMethods PaymentMethods `mapstructure:"payment_methods"`
}
