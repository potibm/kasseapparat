package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	cfgTypes "github.com/potibm/kasseapparat/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	mockCfg := &cfgTypes.Config{
		App: cfgTypes.AppConfig{
			Version:            "1.2.3",
			EnvironmentMessage: "Test-Env",
		},
		Sentry: cfgTypes.SentryConfig{
			DSN: "https://example-dsn.com",
		},
		VATRates: []cfgTypes.VatRateConfig{
			{Rate: 19.0, Name: "Normal"},
		},
		PaymentMethods: []cfgTypes.PaymentMethodConfig{
			{Code: "card", Name: "Kartenzahlung"},
		},
	}

	handler := &Handler{config: *mockCfg}

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetConfig(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Config

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "1.2.3", response.Version)
	assert.Equal(t, "https://example-dsn.com", response.SentryDSN)

	assert.Len(t, response.VATRates, 1)
	assert.Equal(t, 19.0, response.VATRates[0].Rate)

	assert.Len(t, response.PaymentMethods, 1)
	assert.Equal(t, "card", response.PaymentMethods[0].Code)
}
