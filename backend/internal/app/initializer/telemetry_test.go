package initializer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitTelemetry(t *testing.T) {
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
