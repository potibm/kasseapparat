package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"strconv"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := getEnv(key, strconv.Itoa(defaultValue))

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := getEnv(key, strconv.FormatBool(defaultValue))

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	value := getEnv(key, strconv.FormatFloat(defaultValue, 'f', -1, 64))

	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultValue
	}

	return floatValue
}

func getEnvWithJSONValidation(key, fallback string) string {
	val := getEnv(key, fallback)

	var tmp any
	if err := json.Unmarshal([]byte(val), &tmp); err != nil {
		slog.Warn("Invalid JSON", "key", key, "error", err)

		return fallback
	}

	return val
}
