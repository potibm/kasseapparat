package config

import (
	"os"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestReadVersionFromFileWithValidFile(t *testing.T) {
	// Arrange
	filename := "./VERSION"
	expected := "1.2.3\n"
	err := os.WriteFile(filename, []byte(expected), 0644)
	assert.NoError(t, err)

	defer os.Remove(filename) // Clean up

	// Act
	version := readVersionFromFile()

	// Assert
	assert.Equal(t, "1.2.3", version)
}

func TestReadVersionFromFileWithFileMissing(t *testing.T) {
	// Ensure file is absent
	_ = os.Remove("./VERSION")

	version := readVersionFromFile()

	assert.Equal(t, "0.0.0", version)
}

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
			CurrencyLocale: "dk-DK", CurrencyCode: "DKK", DateLocale: "dk-DK", DateOptions: DefaultDateOptions, FractionDigitsMin: 0, FractionDigitsMax: 2,
		},
		VATRates:           DefaultVatRates,
		EnvironmentMessage: "",
		PaymentMethods: PaymentMethods{
			PaymentMethodConfig{
				Name: "💶 Cash", Code: models.PaymentMethodCash},
		},
		SentryConfig: SentryConfig{
			DSN: "", TraceSampleRate: config.SentryConfig.TraceSampleRate, ReplaySessionSampleRate: config.SentryConfig.ReplaySessionSampleRate, ReplayErrorSampleRate: config.SentryConfig.ReplayErrorSampleRate, Environment: "", Version: "0.0.0",
		},
		JwtConfig: JwtConfig{
			Secret: "", Realm: "kasseapparat",
		},
		CorsAllowOrigins: []string{"localhost:3000", "localhost:4000"},
		FrontendURL:      "",
		MailerConfig: MailerConfig{
			DSN: "smtp://user:password@localhost:1025", FromEmail: "", MailSubjectPrefix: "[Kasseapparat]", FrontendURL: "",
		},
		SumupConfig: SumupConfig{
			ApiKey: "", MerchantCode: "", CurrencyCode: "DKK", CurrencyMinorUnit: 2, AffiliateKey: "", ApplicationId: "", PublicUrl: "",
		},
	}

	// Assert
	assert.Equal(t, expected, config)
}
