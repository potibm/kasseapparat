package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitViper(t *testing.T) {
	viper.Reset()

	InitViper()

	// 1. Testing important defaults
	assert.Equal(t, "release", viper.GetString("app.gin_mode"))
	assert.Equal(t, "kasseapparat", viper.GetString("app.db_filename"))
	assert.Equal(t, "da-DK", viper.GetString("format.currency.locale"))
	assert.Equal(t, true, viper.GetBool("jwt.secure_cookie"))

	// 2. Testing the EnvKeyReplacer (Dots to underscores)
	// We set an environment variable simulated via Viper
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Ensure it's loaded
	// If we set APP_LOG_LEVEL, app.log_level should be overridden
	t.Setenv("APP_LOG_LEVEL", "debug")
	assert.Equal(t, "debug", viper.GetString("app.log_level"))

	// 3. Testing the Aliase (That is often the source of errors!)
	viper.Set("app.env", "staging")
	// sentry.environment is an alias for app.env
	assert.Equal(t, "staging", viper.GetString("sentry.environment"))
	assert.Equal(t, "DKK", viper.GetString("sumup.currency_code"))

	viper.Set("app.frontend_url", "https://kasse.party")
	assert.Equal(t, "https://kasse.party", viper.GetString("mailer.frontend_url"))
}

func TestInitViperComplexDefaults(t *testing.T) {
	viper.Reset()
	InitViper()

	vatRates := viper.Get("vatrates")
	assert.NotNil(t, vatRates)

	paymentMethods := viper.Get("payment_methods")
	assert.NotNil(t, paymentMethods)
}
