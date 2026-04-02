package middleware

import (
	"fmt"
	"log/slog"
	"math"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/potibm/kasseapparat/internal/app/models"
)

func payloadFunc() func(data any) jwt.MapClaims {
	return func(data any) jwt.MapClaims {
		foundID := extractUserID(data)

		if foundID != 0 {
			return jwt.MapClaims{
				IdentityKey: foundID,
			}
		}

		slog.Error("JWT payload extraction failed",
			"data_type", fmt.Sprintf("%T", data),
			"content", data,
		)

		return jwt.MapClaims{}
	}
}

func extractUserID(data any) uint {
	switch v := data.(type) {
	case *models.User:
		return v.ID
	case map[string]interface{}:
		return extractIDFromClaims(v)
	default:
		return 0
	}
}

func extractIDFromClaims(claims map[string]interface{}) uint {
	if val, exists := claims["id"]; exists {
		if id, valid := extractUint(val); valid {
			return id
		}
	}

	if val, exists := claims[IdentityKey]; exists {
		if id, valid := extractUint(val); valid {
			return id
		}
	}

	return 0
}

func extractUint(val any) (uint, bool) {
	switch v := val.(type) {
	case float64:
		if v < 0 || v != math.Trunc(v) {
			return 0, false
		}

		return uint(v), true
	case uint:
		return v, true
	case int:
		if v < 0 {
			return 0, false
		}

		return uint(v), true
	case float32:
		f64 := float64(v)
		if f64 < 0 || f64 != math.Trunc(f64) {
			return 0, false
		}

		return uint(v), true
	case int64:
		if v < 0 {
			return 0, false
		}

		return uint(v), true
	default:
		return 0, false // Unsupported type
	}
}
