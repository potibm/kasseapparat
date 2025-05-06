package handler

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentMethodsConfig struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Config struct {
	Version                       string                 `json:"version"`
	SentryDSN                     string                 `json:"sentryDSN"`
	SentryTraceSampleRate         float64                `json:"sentryTraceSampleRate"`
	SentryReplaySessionSampleRate float64                `json:"sentryReplaySessionSampleRate"`
	SentryReplayErrorSampleRate   float64                `json:"sentryReplayErrorSampleRate"`
	CurrencyLocale                string                 `json:"currencyLocale"`
	CurrencyCode                  string                 `json:"currencyCode"`
	VATRates                      string                 `json:"vatRates"`
	DateLocale                    string                 `json:"dateLocale"`
	DateOptions                   string                 `json:"dateOptions"`
	FractionDigitsMin             int                    `json:"fractionDigitsMin"`
	FractionDigitsMax             int                    `json:"fractionDigitsMax"`
	EnvironmentMessage            string                 `json:"environmentMessage"`
	PaymentMethods                []PaymentMethodsConfig `json:"paymentMethods"`
}

func (handler *Handler) GetConfig(c *gin.Context) {
	paymentMethods := make([]PaymentMethodsConfig, 0, len(handler.paymentMethods))
	for code, name := range handler.paymentMethods {
		paymentMethods = append(paymentMethods, PaymentMethodsConfig{
			Code: code,
			Name: name,
		})
	}

	config := Config{
		Version:                       handler.version,
		SentryDSN:                     getEnv("SENTRY_DSN", ""),
		SentryTraceSampleRate:         getEnvAsFloat("SENTRY_TRACE_SAMPLE_RATE", 0.1),
		SentryReplaySessionSampleRate: getEnvAsFloat("SENTRY_REPLAY_SESSION_SAMPLE_RATE", 0.1),
		SentryReplayErrorSampleRate:   getEnvAsFloat("SENTRY_REPLAY_ERROR_SAMPLE_RATE", 0.1),
		CurrencyLocale:                getEnv("CURRENCY_LOCALE", "dk-DK"),
		CurrencyCode:                  getEnv("CURRENCY_CODE", "DKK"),
		VATRates:                      getEnv("VAT_RATES", "[{\"rate\":25,\"name\":\"Standard\"},{\"rate\":0,\"name\":\"Zero rate\"}]"),
		DateLocale:                    getEnv("DATE_LOCALE", "dk-DK"),
		DateOptions:                   getEnv("DATE_OPTIONS", "{\"weekday\":\"long\",\"hour\":\"2-digit\",\"minute\":\"2-digit\"}"),
		FractionDigitsMin:             getEnvAsInt("FRACTION_DIGITS_MIN", 0),
		FractionDigitsMax:             getEnvAsInt("FRACTION_DIGITS_MAX", 2),
		EnvironmentMessage:            getEnv("ENV_MESSAGE", ""),
		PaymentMethods:                paymentMethods,
	}

	c.JSON(200, config)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultValue
	}

	return floatValue
}
