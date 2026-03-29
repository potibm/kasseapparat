package initializer

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"unknown", slog.LevelInfo}, // Fallback
		{"", slog.LevelInfo},        // Leerstring
	}

	for _, tt := range tests {
		t.Run("Level_"+tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, logLevelFromString(tt.input))
		})
	}
}

func TestInitLogger(t *testing.T) {
	t.Run("should set default logger and return it", func(t *testing.T) {
		logger := InitLogger("json", "debug")

		assert.NotNil(t, logger)

		assert.Equal(t, logger, slog.Default())
	})

	t.Run("should handle text type", func(t *testing.T) {
		logger := InitLogger("text", "info")
		assert.NotNil(t, logger)
	})
}
