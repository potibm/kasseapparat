package tests_e2e

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configUrl = "/api/v2/config"

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
	config.Value("fractionDigitsMin").Number().IsEqual(0)
	config.Value("fractionDigitsMax").Number().IsEqual(2)

	dateOptionsRaw := config.Value("dateOptions").String().Raw()

	var dateOptions map[string]any
	if err := json.Unmarshal([]byte(dateOptionsRaw), &dateOptions); err != nil {
		t.Fatalf("dateOptions is not a valid JSON: %v", err)
	}

	for _, key := range []string{"weekday", "hour", "minute"} {
		if _, ok := dateOptions[key]; !ok {
			t.Fatalf("dateOptions does not contain the expected key %q", key)
		}
	}

	assert.Equal(t, "long", dateOptions["weekday"])
	assert.Equal(t, "2-digit", dateOptions["hour"])
	assert.Equal(t, "2-digit", dateOptions["minute"])
}
