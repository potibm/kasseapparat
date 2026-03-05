package initializer

import (
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/potibm/kasseapparat/internal/app/config"
)

func InitializeSentry(sentryConfig config.SentryConfig) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           sentryConfig.DSN,
		Environment:   sentryConfig.Environment,
		EnableTracing: true,
		Release:       "kasseapparat@" + sentryConfig.Version,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: sentryConfig.TraceSampleRate,
	}); err != nil {
		slog.Error("Sentry initialization failed", "error", err)
	}
}
