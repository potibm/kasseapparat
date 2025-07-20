package tests_e2e

import (
	"net/http"
	"testing"
)

var (
	refreshtokenUrl      = "/auth/refresh_token"
	jwtRegexp            = `^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`
	expireDateTimeRegexp = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$`
)

func TestLogin(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	login := e.POST("/login").
		WithJSON(map[string]string{
			"login":    "demo",
			"password": "demo",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	// JWT auslesen
	login.Value("token").String().NotEmpty().Match(jwtRegexp)
	login.Value("code").Number().IsEqual(200)
	login.Value("expire").String().NotEmpty().Match(expireDateTimeRegexp)
	login.Value("role").String().IsEqual("user")
	login.Value("username").String().IsEqual("demo")
	login.Value("id").Number().IsEqual(2)
}

func TestInvalidLogin(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e := e.POST("/login").
		WithJSON(map[string]string{
			"login":    "demo",
			"password": "wrooong",
		}).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object()

	e.Value("code").Number().IsEqual(401)
	e.Value("message").String().IsEqual("incorrect Username or Password")
}

func TestRefreshToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	response := withDemoUserAuthToken(e.GET(refreshtokenUrl)).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	response.Value("code").Number().IsEqual(200)
	response.Value("token").String().NotEmpty().Match(jwtRegexp)
	response.Value("expire").String().NotEmpty().Match(expireDateTimeRegexp)
}

func TestRefreshTokenWithoutToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	response := e.GET(refreshtokenUrl).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object()

	response.Value("code").Number().IsEqual(401)
}

func TestRefreshTokenWithExpiredToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// old token
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNzI2MDg4NTA4LCJvcmlnX2lhdCI6MTcyNjA4NzkwOH0.sNlaDHxoJ6Lr1IruI2DembljhSFmlZncusHUV4hcGq4"

	response := e.GET(refreshtokenUrl).
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object()

	response.Value("code").Number().IsEqual(401)
	response.Value("message").String().IsEqual("Token is expired")
}

func TestRefreshTokenWithInvalidToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// old token
	token := "eythisis.invalid.token"

	response := e.GET(refreshtokenUrl).
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object()

	response.Value("code").Number().IsEqual(401)
	response.Value("message").String().Contains("invalid")
}
