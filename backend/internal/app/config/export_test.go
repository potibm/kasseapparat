package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisUrlUrlObject(t *testing.T) {
	urlStr := "redis://user:password@localhost:6379/0"
	ru := RedisURL(urlStr)

	parsedURL := ru.URLObject()
	assert.NotNil(t, parsedURL)
	assert.Equal(t, "redis", parsedURL.Scheme)
	assert.Equal(t, "localhost:6379", parsedURL.Host)
	assert.Equal(t, "/0", parsedURL.Path)
	assert.Equal(t, "user:password", parsedURL.User.String())

	var nilRu *RedisURL
	assert.Nil(t, nilRu.URLObject())

	invalidURL := RedisURL(":8000") // an url that cant be parsed
	assert.Nil(t, invalidURL.URLObject())
}

func TestRedisUrlIsValid(t *testing.T) {
	validURL := RedisURL("redis://user:password@localhost:6379/0")
	invalidURL := RedisURL("not-a-valid-url")

	assert.True(t, validURL.IsValid())
	assert.False(t, invalidURL.IsValid())
}

func TestRedisUrlJwtConfig(t *testing.T) {
	urlStr := "redis://user:password@localhost:6379/12"
	ru := RedisURL(urlStr)

	jwtConfig := ru.JwtConfig()
	assert.Equal(t, "localhost:6379", jwtConfig.Addr)
	assert.Equal(t, "password", jwtConfig.Password)
	assert.Equal(t, 12, jwtConfig.DB)
}
