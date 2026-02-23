package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	_ = os.Setenv("FOO", "bar")

	assert.Equal(t, "bar", getEnv("FOO", "default"))

	_ = os.Unsetenv("FOO")

	assert.Equal(t, "default", getEnv("FOO", "default"))

	_ = os.Setenv("FOO", "")

	assert.Equal(t, "default", getEnv("FOO", "default"))
}

func TestGetEnvAsInt(t *testing.T) {
	_ = os.Setenv("FOO_INT", "42")

	assert.Equal(t, 42, getEnvAsInt("FOO_INT", 0))

	_ = os.Setenv("FOO_INT", "not-a-number")

	assert.Equal(t, 0, getEnvAsInt("FOO_INT", 0))

	_ = os.Unsetenv("FOO_INT")

	assert.Equal(t, 0, getEnvAsInt("FOO_INT", 0))
}

func TestGetEnvAsFloat(t *testing.T) {
	_ = os.Setenv("FOO_FLOAT", "1.23")

	assert.InDelta(t, 1.23, getEnvAsFloat("FOO_FLOAT", 0), 0.00001)

	_ = os.Setenv("FOO_FLOAT", "not-a-number")

	assert.Equal(t, 3.45, getEnvAsFloat("FOO_FLOAT", 3.45), 0.00001)

	_ = os.Unsetenv("FOO_FLOAT")

	assert.InDelta(t, 2.34, getEnvAsFloat("FOO_FLOAT", 2.34), 0.00001)
}

func TestGetEnvWithJSONValidation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Setenv("VALID_JSON", `[{"x":1}]`)
	t.Setenv("INVALID_JSON", `not-a-json`)
	t.Setenv("EMPTY", "")

	assert.Equal(t, `[{"x":1}]`, getEnvWithJSONValidation(logger, "VALID_JSON", `[{"fallback":true}]`))
	assert.Equal(t, `[{"fallback":true}]`, getEnvWithJSONValidation(logger, "INVALID_JSON", `[{"fallback":true}]`))
	assert.Equal(t, `[{"fallback":true}]`, getEnvWithJSONValidation(logger, "EMPTY", `[{"fallback":true}]`))
}

func TestGetEnvAsBool(t *testing.T) {
	_ = os.Setenv("FOO_BOOL", "true")

	assert.Equal(t, true, getEnvAsBool("FOO_BOOL", false))

	_ = os.Setenv("FOO_BOOL", "false")

	assert.Equal(t, false, getEnvAsBool("FOO_BOOL", true))

	_ = os.Setenv("FOO_BOOL", "not-a-bool")

	assert.Equal(t, false, getEnvAsBool("FOO_BOOL", false))

	_ = os.Unsetenv("FOO_BOOL")

	assert.Equal(t, false, getEnvAsBool("FOO_BOOL", false))
}
