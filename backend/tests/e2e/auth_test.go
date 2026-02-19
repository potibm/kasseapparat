package tests_e2e

import (
	"net/http"
	"testing"
)

var (
	refreshtokenUrl      = "/api/v2/auth/refresh"
	loginUrl             = "/api/v2/auth/login"
	logoutUrl            = "/api/v2/auth/logout"
	jwtRegexp            = `^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`
	expireDateTimeRegexp = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$`
)

func TestLogin(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	r := e.POST(loginUrl).
		WithJSON(map[string]string{
			"login":    "demo",
			"password": "demo",
		}).
		Expect().
		Status(http.StatusOK)

	login := r.JSON().Object()

	login.Value("access_token").String().NotEmpty().Match(jwtRegexp)
	login.Value("expires_in").Number().Gt(0)
	login.Value("role").String().IsEqual("user")
	login.Value("username").String().IsEqual("demo")
	login.Value("id").Number().IsEqual(2)

	refreshToken := r.Cookie("refresh_token")
	refreshToken.Value().NotEmpty()

	jwt := r.Cookie("jwt")
	jwt.Value().NotEmpty()
}

func TestInvalidLogin(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e := e.POST(loginUrl).
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

	req := e.POST(loginUrl).
		WithJSON(map[string]string{
			"login":    "demo",
			"password": "demo",
		}).
		Expect().
		Status(http.StatusOK)

	cookie := req.Cookie("refresh_token")
	oldRefreshToken := cookie.Value().Raw()

	refreshReq := e.POST(refreshtokenUrl).
		WithCookie("refresh_token", oldRefreshToken).
		Expect().
		Status(http.StatusOK)

	refreshResponse := refreshReq.JSON().Object()
	refreshResponse.Value("access_token").String().NotEmpty().Match(jwtRegexp)
	refreshResponse.Value("expires_in").Number().Gt(0)

	cookie2 := refreshReq.Cookie("refresh_token")
	cookie2.Value().NotEmpty().NotEqual(oldRefreshToken)
}

func TestRefreshTokenWithoutToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	response := e.POST(refreshtokenUrl).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object()

	response.Value("code").Number().IsEqual(400)
}

func TestRefreshTokenWithExpiredToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// old token
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
		"eyJJRCI6MiwiZXhwIjoxNzI2MDg4NTA4LCJvcmlnX2lhdCI6MTcyNjA4NzkwOH0." +
		"sNlaDHxoJ6Lr1IruI2DembljhSFmlZncusHUV4hcGq4"

	response := e.POST(refreshtokenUrl).
		WithCookie("refresh_token", token).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object()

	response.Value("code").Number().IsEqual(401)
	response.Value("message").String().IsEqual("invalid or expired refresh token")
}

func TestRefreshTokenWithInvalidToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// old token
	token := "eythisis.invalid.token"

	response := e.POST(refreshtokenUrl).
		WithCookie("refresh_token", token).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object()

	response.Value("code").Number().IsEqual(401)
	response.Value("message").String().Contains("invalid")
}

func TestLogout(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	req := e.POST(loginUrl).
		WithJSON(map[string]string{
			"login":    "demo",
			"password": "demo",
		}).
		Expect().
		Status(http.StatusOK)

	refreshToken := req.Cookie("refresh_token").Value().Raw()
	token := req.Cookie("jwt").Value().Raw()

	logoutReq := e.POST(logoutUrl).
		WithCookie("refresh_token", refreshToken).
		WithCookie("jwt", token).
		WithHeader("Content-Type", "application/json").
		Expect().
		Status(http.StatusOK)

	logoutReq.Cookie("refresh_token").Value().IsEmpty()
	logoutReq.Cookie("jwt").Value().IsEmpty()

	e.POST(refreshtokenUrl).
		WithCookie("refresh_token", refreshToken).
		Expect().
		Status(http.StatusUnauthorized)
}
