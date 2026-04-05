package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVatRatesConfigJson(t *testing.T) {
	vr := VatRatesConfig{
		{Name: "Standard", Rate: DefaultStandardVatRate},
	}

	expectedJson := `[{"Rate":25,"Name":"Standard"}]`
	assert.Equal(t, expectedJson, vr.Json())
}

func TestDateFormatOptionsConfigJson(t *testing.T) {
	dfo := DateFormatOptionsConfig{
		"weekday": "long",
	}

	expectedJson := `{"weekday":"long"}`
	assert.Equal(t, expectedJson, dfo.Json())
}

func TestRedisUrlUrlObject(t *testing.T) {
	urlStr := "redis://user:password@localhost:6379/0"
	ru := RedisUrl(urlStr)

	parsedUrl := ru.UrlObject()
	assert.NotNil(t, parsedUrl)
	assert.Equal(t, "redis", parsedUrl.Scheme)
	assert.Equal(t, "localhost:6379", parsedUrl.Host)
	assert.Equal(t, "/0", parsedUrl.Path)
	assert.Equal(t, "user:password", parsedUrl.User.String())

	var nilRu *RedisUrl
	assert.Nil(t, nilRu.UrlObject())

	invalidUrl := RedisUrl(":8000") // an url that cant be parsed
	assert.Nil(t, invalidUrl.UrlObject())
}

func TestRedisUrlIsValid(t *testing.T) {
	validUrl := RedisUrl("redis://user:password@localhost:6379/0")
	invalidUrl := RedisUrl("not-a-valid-url")

	assert.True(t, validUrl.IsValid())
	assert.False(t, invalidUrl.IsValid())
}

func TestRedisUrlJwtConfig(t *testing.T) {
	urlStr := "redis://user:password@localhost:6379/12"
	ru := RedisUrl(urlStr)

	jwtConfig := ru.JwtConfig()
	assert.Equal(t, "localhost:6379", jwtConfig.Addr)
	assert.Equal(t, "password", jwtConfig.Password)
	assert.Equal(t, 12, jwtConfig.DB)
}
