package config

import (
	"fmt"

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

type AppConfig struct {
	Version string `mapstructure:"version"`

	GinMode     string `mapstructure:"gin_mode" validate:"required,oneof=debug release test"`
	Environment string `mapstructure:"env"      validate:"required,oneof=development staging production test"`

	LogLevel  string `mapstructure:"log_level"  validate:"required,oneof=debug info warn error"`
	LogFormat string `mapstructure:"log_format" validate:"required,oneof=json text"`

	DbFilename         string                 `mapstructure:"db_filename"         validate:"required"`
	RedisURL           RedisURL               `mapstructure:"redis_url"           validate:"omitempty,url"`
	FrontendURL        string                 `mapstructure:"frontend_url"        validate:"required,http_url"`
	CorsAllowOrigins   CorsAllowOriginsConfig `mapstructure:"cors_allow_origins"  validate:"dive,required"`
	EnvironmentMessage string                 `mapstructure:"environment_message"`

	OtelEndpoint string `mapstructure:"otel_endpoint" validate:"omitempty"`
	Port         int    `mapstructure:"port"          validate:"required,gt=0,lte=65535"`
}

type FormatConfig struct {
	Currency CurrencyFormatConfig `mapstructure:"currency"`
	Date     DateFormatConfig     `mapstructure:"date"`
}

type CurrencyFormatConfig struct {
	Locale            string `mapstructure:"locale"              validate:"required"`
	Code              string `mapstructure:"code"                validate:"required"`
	FractionDigitsMin int32  `mapstructure:"fraction_digits_min" validate:"gte=0"`
	FractionDigitsMax int32  `mapstructure:"fraction_digits_max" validate:"gte=0"`
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
	APIKey            string `mapstructure:"api_key"`
	MerchantCode      string `mapstructure:"merchant_code"`
	CurrencyCode      string `mapstructure:"currency_code"`
	CurrencyMinorUnit int32  `mapstructure:"currency_minor_unit"`
	AffiliateKey      string `mapstructure:"affiliate_key"`
	ApplicationID     string `mapstructure:"application_id"`
	PublicURL         string `mapstructure:"public_url"          validate:"omitempty,https_url"`
}

type AuthConfig struct {
	Mode        string   `mapstructure:"mode"         validate:"required,oneof=proxy oidc"`
	ProxyHeader string   `mapstructure:"proxy_header"`
	ProxyAdmins []string `mapstructure:"proxy_admins"`

	OidcIssuer       string `mapstructure:"oidc_issuer"`
	OidcClientID     string `mapstructure:"oidc_client_id"`
	OidcClientSecret string `mapstructure:"oidc_client_secret"`
	OidcCallbackURL  string `mapstructure:"oidc_callback_url"`
	SessionSecret    string `mapstructure:"session_secret"`
}

func (a *AuthConfig) Validate() error {
	if a.Mode == "proxy" {
		if a.ProxyHeader == "" {
			return fmt.Errorf("auth.proxy_header is required when mode is 'proxy'")
		}

		return nil
	}

	if a.Mode == "oidc" {
		return a.validateOIDC()
	}

	return nil
}

const MinSessionSecretLength = 32

func (a *AuthConfig) validateOIDC() error {
	if a.OidcIssuer == "" {
		return fmt.Errorf("auth.oidc_issuer is required when mode is 'oidc'")
	}

	if a.OidcClientID == "" {
		return fmt.Errorf("auth.oidc_client_id is required when mode is 'oidc'")
	}

	if a.OidcClientSecret == "" {
		return fmt.Errorf("auth.oidc_client_secret is required when mode is 'oidc'")
	}

	if a.OidcCallbackURL == "" {
		return fmt.Errorf("auth.oidc_callback_url is required when mode is 'oidc'")
	}

	if len(a.SessionSecret) < MinSessionSecretLength {
		return fmt.Errorf(
			"auth.session_secret must be at least %d characters when mode is 'oidc'",
			MinSessionSecretLength,
		)
	}

	return nil
}

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Format FormatConfig `mapstructure:"format"`
	Sentry SentryConfig `mapstructure:"sentry"`
	Jwt    JwtConfig    `mapstructure:"jwt"`
	Mailer MailerConfig `mapstructure:"mailer"`
	Sumup  SumupConfig  `mapstructure:"sumup"`
	Auth   AuthConfig   `mapstructure:"auth"`

	VATRates       VatRatesConfig `mapstructure:"vat_rates"`
	PaymentMethods PaymentMethods `mapstructure:"payment_methods"`
}
