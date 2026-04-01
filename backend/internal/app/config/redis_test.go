package config

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validRedisURL  = "redis://:password@localhost:6379/0"
	validRedisHost = "localhost:6379"
)

func TestLoadRedisConfigWithEmptyEnv(t *testing.T) {
	os.Setenv("REDIS_URL", "")

	defer os.Unsetenv("REDIS_URL")

	cfg, err := loadRedisConfig()

	assert.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestLoadRedisConfigWithValidEnv(t *testing.T) {
	os.Setenv("REDIS_URL", validRedisURL)

	defer os.Unsetenv("REDIS_URL")

	cfg, err := loadRedisConfig()

	assert.NoError(t, err)

	if !assert.NotNil(t, cfg) {
		return
	}

	assert.Equal(t, validRedisHost, cfg.Host)
	assert.Equal(t, "/0", cfg.Path)

	if assert.NotNil(t, cfg.User) {
		password, isSet := cfg.User.Password()
		assert.True(t, isSet, "password should be set")
		assert.Equal(t, "password", password)
	}
}

func TestLoadRedisConfigWithInvalidUrlEnv(t *testing.T) {
	os.Setenv("REDIS_URL", "redis://localhost/%zz")

	defer os.Unsetenv("REDIS_URL")

	cfg, err := loadRedisConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "invalid URL escape")
}

func TestLoadRedisConfigWithInvalidSchemeEnv(t *testing.T) {
	os.Setenv("REDIS_URL", "https://localhost:6379/0")

	defer os.Unsetenv("REDIS_URL")

	cfg, err := loadRedisConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "invalid scheme")
}

func TestJwtConfigWithValidRedisConfig(t *testing.T) {
	redisConfig := RedisConfig{
		Scheme: "redis",
		User:   url.UserPassword("username", "password"),
		Host:   validRedisHost,
		Path:   "/1",
	}

	jwtConfig := redisConfig.JwtConfig()

	assert.Equal(t, validRedisHost, jwtConfig.Addr)
	assert.Equal(t, "password", jwtConfig.Password)
	assert.Equal(t, 1, jwtConfig.DB)
}

func TestJwtConfigWithMissingPath(t *testing.T) {
	redisConfig := RedisConfig{
		Scheme: "redis",
		User:   url.UserPassword("username", "password"),
		Host:   validRedisHost,
	}

	jwtConfig := redisConfig.JwtConfig()

	assert.Equal(t, validRedisHost, jwtConfig.Addr)
	assert.Equal(t, "password", jwtConfig.Password)
	assert.Equal(t, 0, jwtConfig.DB)
}

func TestJwtConfigWithInvalidPath(t *testing.T) {
	redisConfig := RedisConfig{
		Scheme: "redis",
		User:   url.UserPassword("username", "password"),
		Host:   validRedisHost,
		Path:   "/invalid",
	}

	jwtConfig := redisConfig.JwtConfig()

	assert.Equal(t, validRedisHost, jwtConfig.Addr)
	assert.Equal(t, "password", jwtConfig.Password)
	assert.Equal(t, 0, jwtConfig.DB)
}

func TestJwtConfigWithoutPassword(t *testing.T) {
	redisConfig := RedisConfig{
		Scheme: "redis",
		User:   url.User("username"),
		Host:   validRedisHost,
		Path:   "/0",
	}

	jwtConfig := redisConfig.JwtConfig()

	assert.Equal(t, validRedisHost, jwtConfig.Addr)
	assert.Equal(t, "", jwtConfig.Password)
}
