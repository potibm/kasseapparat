package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisUrlJwtConfig(t *testing.T) {
	urlStr := "redis://user:password@localhost:6379/12"
	ru := RedisURL(urlStr)

	jwtConfig := ru.JwtConfig()
	assert.Equal(t, "localhost:6379", jwtConfig.Addr)
	assert.Equal(t, "password", jwtConfig.Password)
	assert.Equal(t, 12, jwtConfig.DB)
}
