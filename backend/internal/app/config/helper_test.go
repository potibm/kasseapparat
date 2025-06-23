package config

import (
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
	t.Setenv("VALID_JSON", `[{"x":1}]`)
	t.Setenv("INVALID_JSON", `not-a-json`)
	t.Setenv("EMPTY", "")

	assert.Equal(t, `[{"x":1}]`, getEnvWithJSONValidation("VALID_JSON", `[{"fallback":true}]`))
	assert.Equal(t, `[{"fallback":true}]`, getEnvWithJSONValidation("INVALID_JSON", `[{"fallback":true}]`))
	assert.Equal(t, `[{"fallback":true}]`, getEnvWithJSONValidation("EMPTY", `[{"fallback":true}]`))
}
