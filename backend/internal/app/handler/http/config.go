package http

import (
	nethttp "net/http"

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
	paymentMethods := make([]PaymentMethodsConfig, 0, len(handler.config.PaymentMethods))

	for _, configPaymentMethod := range handler.config.PaymentMethods {
		paymentMethods = append(paymentMethods, PaymentMethodsConfig{
			Code: string(configPaymentMethod.Code),
			Name: configPaymentMethod.Name,
		})
	}

	config := Config{
		Version:                       handler.config.App.Version,
		SentryDSN:                     handler.config.Sentry.DSN,
		SentryTraceSampleRate:         handler.config.Sentry.TraceSampleRate,
		SentryReplaySessionSampleRate: handler.config.Sentry.ReplaySessionSampleRate,
		SentryReplayErrorSampleRate:   handler.config.Sentry.ReplayErrorSampleRate,
		CurrencyLocale:                handler.config.Format.Currency.Locale,
		CurrencyCode:                  handler.config.Format.Currency.Code,
		FractionDigitsMin:             handler.config.Format.Currency.FractionDigitsMin,
		FractionDigitsMax:             handler.config.Format.Currency.FractionDigitsMax,
		VATRates:                      handler.config.VATRates.Json(),
		DateLocale:                    handler.config.Format.Date.Locale,
		DateOptions:                   handler.config.Format.Date.Options.Json(),
		EnvironmentMessage:            handler.config.App.EnvironmentMessage,
		PaymentMethods:                paymentMethods,
	}

	c.JSON(nethttp.StatusOK, config)
}
