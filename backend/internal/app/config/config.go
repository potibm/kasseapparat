package config

import (
	"strings"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/spf13/viper"
)

const (
	OtelServiceName = "kasseapparat-backend"

	DefaultTraceSampleRate         = 0.1
	DefaultReplaySessionSampleRate = 0.1
	DefaultReplayErrorSampleRate   = 0.1
	DefaultMinorUnit               = 2
	DefaultJwtSecret               = "very-insecure"

	DefaultStandardVatRate = 25
	DefaultReducedVatRate  = 12
	DefaultZeroVatRate     = 0
)

var (
	DefaultDateOptions = DateFormatOptionsConfig{
		"weekday": "long",
		"hour":    "2-digit",
		"minute":  "2-digit",
	}
	DefaultVatRates = VatRatesConfig{
		{Name: "Standard", Rate: DefaultStandardVatRate},
		{Name: "Reduced", Rate: DefaultReducedVatRate},
		{Name: "Zero", Rate: DefaultZeroVatRate},
	}
	DefaultPaymentMethods = []PaymentMethodConfig{
		{Code: models.PaymentMethodCash, Name: "Cash"},
		{Code: models.PaymentMethodCC, Name: "Creditcard"},
		{Code: models.PaymentMethodVoucher, Name: "Voucher"},
	}
)

func InitViper() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("app.gin_mode", "release")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.env", "production")
	viper.SetDefault("app.db_filename", "kasseapparat")
	viper.SetDefault("app.redis_url", "")
	viper.SetDefault("app.frontend_url", "")
	viper.SetDefault("app.cors_allow_origins", []string{})

	viper.SetDefault("format.currency.locale", "da-DK")
	viper.SetDefault("format.currency.code", "DKK")
	viper.SetDefault("format.currency.fraction_digits_min", 0)
	viper.SetDefault("format.currency.fraction_digits_max", DefaultMinorUnit)
	viper.SetDefault("format.date.locale", "da-DK")
	viper.SetDefault("format.date.options", DefaultDateOptions)

	viper.SetDefault("sentry.dsn", "")
	viper.SetDefault("sentry.trace_sample_rate", DefaultTraceSampleRate)
	viper.SetDefault("sentry.replay_session_sample_rate", DefaultReplaySessionSampleRate)
	viper.SetDefault("sentry.replay_error_sample_rate", DefaultReplayErrorSampleRate)

	viper.SetDefault("jwt.realm", "Kasseapparat")
	viper.SetDefault("jwt.secret", DefaultJwtSecret)
	viper.SetDefault("jwt.secure_cookie", true)

	viper.SetDefault("mailer.dsn", "smtp://user:password@localhost:1025")
	viper.SetDefault("mailer.from", "kasseapparat@example.com")
	viper.SetDefault("mailer.subject_prefix", "[Kasseapparat]")

	viper.SetDefault("sumup.api_key", "")
	viper.SetDefault("sumup.merchant_code", "")
	viper.SetDefault("sumup.currency_code", "")
	viper.SetDefault("sumup.currency_minor_unit", DefaultMinorUnit)
	viper.SetDefault("sumup.affiliate_key", "")
	viper.SetDefault("sumup.application_id", "")
	viper.SetDefault("sumup.public_url", "")

	viper.SetDefault("vatrates", DefaultVatRates)
	viper.SetDefault("payment_methods", DefaultPaymentMethods)

	viper.RegisterAlias("mailer.frontend_url", "app.frontend_url")
	viper.RegisterAlias("sentry.environment", "app.env")
	viper.RegisterAlias("sentry.version", "app.version")
	viper.RegisterAlias("sumup.currency_code", "format.currency.code")
}
