package initializer

import (
	"fmt"
	"os"
	"strconv"

	"github.com/getsentry/sentry-go"
)

func InitializeSentry() {
	tracesSampleRate, err := strconv.ParseFloat(os.Getenv("SENTRY_TRACE_SAMPLE_RATE"), 32);
	if err != nil {
		tracesSampleRate = 0.1
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		Environment: os.Getenv("SENTRY_ENVIRONMENT"),
		EnableTracing: true,
		Release: "kasseapparat@" + GetVersion(),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: tracesSampleRate,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
}
