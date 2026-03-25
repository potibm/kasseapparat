package initializer

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/potibm/kasseapparat/internal/app/config"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
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

func InitLogger(ctx context.Context, loggerType, level string) *slog.Logger {
	if loggerType == "text" {
		return InitTxtLogger(ctx, level)
	}

	return InitJsonLogger(ctx, level) // default to JSON logger
}

func InitJsonLogger(ctx context.Context, level string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevelFromString(level),
	})

	return initializeLogger(ctx, handler)
}

func InitTxtLogger(ctx context.Context, level string) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevelFromString(level),
	})

	return initializeLogger(ctx, handler)
}

func initializeLogger(ctx context.Context, cmdlineHandler slog.Handler) *slog.Logger {
	otelHandler := otelslog.NewHandler(config.OtelServiceName);

	var finalHandler slog.Handler
	
	finalHandler = slogmulti.Fanout(cmdlineHandler, otelHandler)
	
	logger := slog.New(finalHandler)
	slog.SetDefault(logger)

	return logger
}
