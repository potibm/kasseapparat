package initializer

import (
	"log/slog"
	"os"
	"strings"
)

func logLevelFromString(level string) slog.Level {
	level = strings.ToLower(level)

	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // default to info if unrecognized
	}
}

func InitLogger(loggerType, level string) *slog.Logger {
	if loggerType == "text" {
		return InitTxtLogger(level)
	}

	return InitJsonLogger(level) // default to JSON logger
}

func InitJsonLogger(level string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevelFromString(level),
	})

	return initializeLogger(handler)
}

func InitTxtLogger(level string) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevelFromString(level),
	})

	return initializeLogger(handler)
}

func initializeLogger(handler slog.Handler) *slog.Logger {
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
