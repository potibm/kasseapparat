package middleware

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestExtractUint(t *testing.T) {
	tests := []struct {
		name          string
		input         any
		expectedID    uint
		expectedValid bool
	}{
		{"Float64 from JSON", float64(123), 123, true},
		{"Float32", float32(123), 123, true},
		{"Native int", int(789), 789, true},
		{"Native int64", int64(789), 789, true},
		{"Native uint", uint(456), 456, true},
		{"Native uint64", uint(456), 456, true},
		{"Negative int", int(-789), 0, false},
		{"Negative Float64 ", float64(-123), 0, false},
		{"Negative Float32", float32(-234), 0, false},
		{"Negative int64", int64(-345), 0, false},
		{"Negative int", int(-456), 0, false},
		{"Fractional Float64", float64(123.456), 0, false},
		{"Fractional Float32", float32(123.456), 0, false},
		{"Invalid string type", "123", 0, false},
		{"Nil value", nil, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, valid := extractUint(tt.input)
			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedValid, valid)
		})
	}
}

func TestExtractIDFromClaims(t *testing.T) {
	tests := []struct {
		name       string
		claims     map[string]interface{}
		expectedID uint
	}{
		{"Lowercase id exists", map[string]interface{}{"id": float64(111)}, 111},
		{"IdentityKey exists", map[string]interface{}{IdentityKey: float64(222)}, 222},
		{"Lowercase id has priority", map[string]interface{}{"id": float64(111), IdentityKey: float64(222)}, 111},
		{"No valid key exists", map[string]interface{}{"username": "demo"}, 0},
		{"Key exists but wrong type", map[string]interface{}{"id": "not-a-number"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := extractIDFromClaims(tt.claims)
			assert.Equal(t, tt.expectedID, id)
		})
	}
}

func TestExtractUserID(t *testing.T) {
	t.Run("Valid User Struct", func(t *testing.T) {
		user := &models.User{ID: 999}
		id := extractUserID(user)
		assert.Equal(t, uint(999), id)
	})

	t.Run("Valid Claims Map", func(t *testing.T) {
		claims := map[string]interface{}{"id": float64(888)}
		id := extractUserID(claims)
		assert.Equal(t, uint(888), id)
	})

	t.Run("Invalid Type", func(t *testing.T) {
		id := extractUserID("just a string")
		assert.Equal(t, uint(0), id)
	})
}

func TestPayloadFunc(t *testing.T) {
	t.Run("Success with User struct", func(t *testing.T) {
		f := payloadFunc()
		user := &models.User{ID: 123}
		claims := f(user)

		assert.Equal(t, uint(123), claims[IdentityKey])
	})

	t.Run("Failure triggers slog Error and returns empty claims", func(t *testing.T) {
		var buf bytes.Buffer

		testLogger := slog.New(slog.NewTextHandler(&buf, nil))
		oldLogger := slog.Default()

		slog.SetDefault(testLogger)
		defer slog.SetDefault(oldLogger)

		f := payloadFunc()
		claims := f("invalid data type")

		// Assertions
		assert.Empty(t, claims)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "JWT payload extraction failed")
		assert.Contains(t, logOutput, "level=ERROR")
	})
}
