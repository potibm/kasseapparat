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

func InitJsonLogger(level string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevelFromString(level),
	})

	return initLogger(handler)
}

func InitTxtLogger(level string) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevelFromString(level),
	})

	return initLogger(handler)
}

func initLogger(handler slog.Handler) *slog.Logger {
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
