package tests_e2e

import (
	"net/http"
	"testing"
)

var (
	configUrl = "/api/v1/config"
)

func TestGetConfig(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	config := e.GET(configUrl).
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
	config.Value("dateOptions").String().Contains("{\"weekday\":\"long\",\"hour\":\"2-digit\",\"minute\":\"2-digit\"}")
	config.Value("fractionDigitsMin").Number().IsEqual(0)
	config.Value("fractionDigitsMax").Number().IsEqual(2)
}
