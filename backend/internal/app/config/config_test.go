package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigWithDefaults(t *testing.T) {
	os.Setenv("CORS_ALLOW_ORIGINS", "localhost:3000,localhost:4000")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Act
	config, err := loadConfig(logger)

	assert.NoError(t, err)

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
			TraceSampleRate:         defaultTraceSampleRate,
			ReplaySessionSampleRate: defaultReplaySessionSampleRate,
			ReplayErrorSampleRate:   defaultReplayErrorSampleRate,
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

	assert.Equal(t, expected.EnvironmentMessage, config.EnvironmentMessage)
	assert.Equal(t, expected.FrontendURL, config.FrontendURL)

	// test all the other fields in a loop to avoid writing a lot of repetitive code
	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
	}{
		{"AppConfig", expected.AppConfig, config.AppConfig},
		{"FormatConfig", expected.FormatConfig, config.FormatConfig},
		{"VATRates", expected.VATRates, config.VATRates},
		{"PaymentMethods", expected.PaymentMethods, config.PaymentMethods},
		{"JwtConfig", expected.JwtConfig, config.JwtConfig},
		{"CorsAllowOrigins", expected.CorsAllowOrigins, config.CorsAllowOrigins},
		{"MailerConfig", expected.MailerConfig, config.MailerConfig},
		{"SumupConfig", expected.SumupConfig, config.SumupConfig},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.actual)
			assert.Equal(t, tt.expected, tt.actual)
		})
	}

	// sentry config needs to be tested separately because of the float fields
	assert.NotNil(t, config.SentryConfig)
	assert.Equal(t, expected.SentryConfig.DSN, config.SentryConfig.DSN)
	assert.Equal(t, expected.SentryConfig.Environment, config.SentryConfig.Environment)
	assert.Equal(t, expected.SentryConfig.Version, config.SentryConfig.Version)
	assert.InEpsilon(t, expected.SentryConfig.TraceSampleRate, config.SentryConfig.TraceSampleRate, 0.0001)
	assert.InEpsilon(
		t,
		expected.SentryConfig.ReplaySessionSampleRate,
		config.SentryConfig.ReplaySessionSampleRate,
		0.0001,
	)
	assert.InEpsilon(t, expected.SentryConfig.ReplayErrorSampleRate, config.SentryConfig.ReplayErrorSampleRate, 0.0001)
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
