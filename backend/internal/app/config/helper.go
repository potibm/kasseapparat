package config

import (
	"encoding/json"
	"log"
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
		log.Printf("Invalid JSON in %s: %v", key, err)

		return fallback
	}

	return val
}
