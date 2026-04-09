package tests_e2e

import (
	"net/http"
	"testing"
)

var configURL = "/api/v2/config"

func TestGetConfig(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	config := e.GET(configURL).
		Expect().
		Status(http.StatusOK).JSON().Object()

	config.Value("version").String()
	config.Value("sentryDSN").String()
	config.Value("sentryTraceSampleRate").Number()
	config.Value("sentryReplaySessionSampleRate").Number()
	config.Value("sentryReplayErrorSampleRate").Number()
	config.Value("currencyLocale").String().Match("^[a-z]{2}-[A-Z]{2}$")
	config.Value("currencyCode").String().Match("^[A-Z]{3}$")
	config.Value("dateLocale").String().Match("^[a-z]{2}-[A-Z]{2}$")
	config.Value("fractionDigitsMin").Number().IsEqual(0)
	config.Value("fractionDigitsMax").Number().IsEqual(2)

	config.Value("dateOptions").Object()
	config.Value("dateOptions").Object().Value("weekday").IsEqual("long")
	config.Value("dateOptions").Object().Value("hour").IsEqual("2-digit")
	config.Value("dateOptions").Object().Value("minute").IsEqual("2-digit")

	paymentMethods := config.Value("paymentMethods").Array()
	paymentMethods.NotEmpty() // Ersetzt assert.Greater(..., 0)

	for _, item := range paymentMethods.Iter() {
		obj := item.Object()
		obj.Value("code").String().NotEmpty()
		obj.Value("name").String().NotEmpty()
	}

	vatRates := config.Value("vatRates").Array()
	vatRates.NotEmpty()

	for _, item := range vatRates.Iter() {
		obj := item.Object()
		obj.Value("rate").Number()
		obj.Value("name").String().NotEmpty()
	}
}
