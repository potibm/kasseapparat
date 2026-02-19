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
		Version:                       handler.config.AppConfig.Version,
		SentryDSN:                     handler.config.SentryConfig.DSN,
		SentryTraceSampleRate:         handler.config.SentryConfig.TraceSampleRate,
		SentryReplaySessionSampleRate: handler.config.SentryConfig.ReplaySessionSampleRate,
		SentryReplayErrorSampleRate:   handler.config.SentryConfig.ReplayErrorSampleRate,
		CurrencyLocale:                handler.config.FormatConfig.CurrencyLocale,
		CurrencyCode:                  handler.config.FormatConfig.CurrencyCode,
		VATRates:                      handler.config.VATRates,
		DateLocale:                    handler.config.FormatConfig.DateLocale,
		DateOptions:                   handler.config.FormatConfig.DateOptions,
		FractionDigitsMin:             handler.config.FormatConfig.FractionDigitsMin,
		FractionDigitsMax:             handler.config.FormatConfig.FractionDigitsMax,
		EnvironmentMessage:            handler.config.EnvironmentMessage,
		PaymentMethods:                paymentMethods,
	}

	c.JSON(nethttp.StatusOK, config)
}
