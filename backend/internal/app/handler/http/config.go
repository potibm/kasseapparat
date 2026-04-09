package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	cfgTypes "github.com/potibm/kasseapparat/internal/app/config"
)

type PaymentMethodsConfig struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type VatRateConfig struct {
	Rate float64 `json:"rate"`
	Name string  `json:"name"`
}

type DateFormatOptionsConfig map[string]any

type Config struct {
	Version                       string                  `json:"version"`
	SentryDSN                     string                  `json:"sentryDSN"`
	SentryTraceSampleRate         float64                 `json:"sentryTraceSampleRate"`
	SentryReplaySessionSampleRate float64                 `json:"sentryReplaySessionSampleRate"`
	SentryReplayErrorSampleRate   float64                 `json:"sentryReplayErrorSampleRate"`
	CurrencyLocale                string                  `json:"currencyLocale"`
	CurrencyCode                  string                  `json:"currencyCode"`
	VATRates                      []VatRateConfig         `json:"vatRates"`
	DateLocale                    string                  `json:"dateLocale"`
	DateOptions                   DateFormatOptionsConfig `json:"dateOptions"`
	FractionDigitsMin             int32                   `json:"fractionDigitsMin"`
	FractionDigitsMax             int32                   `json:"fractionDigitsMax"`
	EnvironmentMessage            string                  `json:"environmentMessage"`
	PaymentMethods                []PaymentMethodsConfig  `json:"paymentMethods"`
}

func (handler *Handler) GetConfig(c *gin.Context) {
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
		VATRates:                      convertVatRates(handler.config.VATRates),
		DateLocale:                    handler.config.Format.Date.Locale,
		DateOptions:                   DateFormatOptionsConfig(handler.config.Format.Date.Options),
		EnvironmentMessage:            handler.config.App.EnvironmentMessage,
		PaymentMethods:                convertPaymentMethods(handler.config.PaymentMethods),
	}

	c.JSON(nethttp.StatusOK, config)
}

func convertPaymentMethods(paymentMethods []cfgTypes.PaymentMethodConfig) []PaymentMethodsConfig {
	result := make([]PaymentMethodsConfig, 0, len(paymentMethods))

	for _, configPaymentMethod := range paymentMethods {
		result = append(result, PaymentMethodsConfig{
			Code: string(configPaymentMethod.Code),
			Name: configPaymentMethod.Name,
		})
	}

	return result
}

func convertVatRates(vatRates []cfgTypes.VatRateConfig) []VatRateConfig {
	result := make([]VatRateConfig, 0, len(vatRates))

	for _, configVatRate := range vatRates {
		result = append(result, VatRateConfig{
			Rate: configVatRate.Rate,
			Name: configVatRate.Name,
		})
	}

	return result
}
