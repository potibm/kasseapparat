package initializer

import (
	"errors"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/stretchr/testify/assert"
)

const (
	jwtTestReal = "test-realm"
)

func TestInitializeJwtMiddlewareWithoutRedis(t *testing.T) {
	// 1. Arrange
	jwtConfig := config.JwtConfig{
		Realm:        jwtTestReal,
		Secret:       "super-secret-key-that-is-long-enough",
		SecureCookie: true,
	}
	// 2. Act
	middleware := InitializeJwtMiddleware(nil, jwtConfig, nil)

	// 3. Assert
	assert.NotNil(t, middleware)
	assert.Equal(t, jwtTestReal, middleware.Realm)
}

func TestInitializeJwtMiddlewareWithRedis(t *testing.T) {
	// 1. Arrange
	jwtConfig := config.JwtConfig{
		Realm:        jwtTestReal,
		Secret:       "super-secret-key-that-is-long-enough",
		SecureCookie: true,
	}

	redisURL := config.RedisURL("redis://localhost:6379/0")

	// 2. Act
	middleware := InitializeJwtMiddleware(nil, jwtConfig, &redisURL)

	// 3. Assert
	assert.NotNil(t, middleware)
}

func TestInitializeJwtMiddlewarePanicOnError(t *testing.T) {
	originalJwtFunc := newJwtFunc

	defer func() { newJwtFunc = originalJwtFunc }()

	newJwtFunc = func(m *jwt.GinJWTMiddleware) (*jwt.GinJWTMiddleware, error) {
		return nil, errors.New("forced error for test coverage")
	}

	jwtConfig := config.JwtConfig{
		Realm:        jwtTestReal,
		Secret:       "",
		SecureCookie: true,
	}

	assert.Panics(t, func() {
		InitializeJwtMiddleware(nil, jwtConfig, nil)
	}, "This should panic due to invalid JWT configuration")
}
