package initializer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
)

func TestInitTelemetry(t *testing.T) {
	oldTP := otel.GetTracerProvider()
	oldMP := otel.GetMeterProvider()
	oldLP := global.GetLoggerProvider()

	t.Cleanup(func() {
		otel.SetTracerProvider(oldTP)
		otel.SetMeterProvider(oldMP)
		global.SetLoggerProvider(oldLP)
	})

	ctx := context.Background()
	version := "1.0.0"

	t.Run("should return nil if endpoint is empty", func(t *testing.T) {
		cleanup, err := InitTelemetry(ctx, "", version)
		assert.NoError(t, err)
		assert.Nil(t, cleanup)
	})

	t.Run("should initialize with insecure endpoint", func(t *testing.T) {
		endpoint := "localhost:4317"

		cleanup, err := InitTelemetry(ctx, endpoint, version)

		require.NoError(t, err)
		require.NotNil(t, cleanup)

		// run cleanup function to ensure it doesn't panic
		assert.NotPanics(t, func() {
			cleanup()
		})
	})
}
