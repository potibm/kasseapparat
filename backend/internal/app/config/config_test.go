package config

import (
	"os"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigWithDefaults(t *testing.T) {
	os.Setenv("CORS_ALLOW_ORIGINS", "localhost:3000,localhost:4000")

	// Act
	config := loadConfig()

	// Arrange
	expected := Config{
		AppConfig: AppConfig{
			Version: "0.0.0",
			GinMode: "release",
		},
		FormatConfig: FormatConfig{
			CurrencyLocale:    "dk-DK",
			CurrencyCode:      "DKK",
			DateLocale:        "dk-DK",
			DateOptions:       DefaultDateOptions,
			FractionDigitsMin: 0,
			FractionDigitsMax: 2,
		},
		VATRates:           DefaultVatRates,
		EnvironmentMessage: "",
		PaymentMethods: PaymentMethods{
			PaymentMethodConfig{
				Name: "💶 Cash", Code: models.PaymentMethodCash},
		},
		SentryConfig: SentryConfig{
			DSN:                     "",
			TraceSampleRate:         config.SentryConfig.TraceSampleRate,
			ReplaySessionSampleRate: config.SentryConfig.ReplaySessionSampleRate,
			ReplayErrorSampleRate:   config.SentryConfig.ReplayErrorSampleRate,
			Environment:             "",
			Version:                 "0.0.0",
		},
		JwtConfig: JwtConfig{
			Secret: "", Realm: "kasseapparat", SecureCookie: true,
		},
		CorsAllowOrigins: []string{"localhost:3000", "localhost:4000"},
		FrontendURL:      "",
		MailerConfig: MailerConfig{
			DSN:               "smtp://user:password@localhost:1025",
			FromEmail:         "",
			MailSubjectPrefix: "[Kasseapparat]",
			FrontendURL:       "",
		},
		SumupConfig: SumupConfig{
			ApiKey:            "",
			MerchantCode:      "",
			CurrencyCode:      "DKK",
			CurrencyMinorUnit: 2,
			AffiliateKey:      "",
			ApplicationId:     "",
			PublicUrl:         "",
		},
	}

	// Assert
	assert.Equal(t, expected, config)
}

func TestSetVersion(t *testing.T) {
	// Arrange
	config := Config{
		AppConfig: AppConfig{
			Version: "0.0.0",
		},
		SentryConfig: SentryConfig{
			Version: "0.0.0",
		},
	}

	config.SetVersion("1.2.3")

	assert.Equal(t, "1.2.3", config.AppConfig.Version)
	assert.Equal(t, "1.2.3", config.SentryConfig.Version)
}
