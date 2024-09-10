package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"
)

func TestGetUsers(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := e.GET("/api/v1/users").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)

	totalCountHeader := res.Header("X-Total-Count").AsNumber()
	totalCountHeader.Ge(2)

	obj := res.JSON().Array()
	obj.Length().Ge(2)

	for i := 0; i < len(obj.Iter()); i++ {
		user := obj.Value(i).Object()
		user.Value("id").Number().Gt(0)
		user.Value("username").String().NotEmpty()
		user.Value("email").String().NotEmpty().Contains("@")
		user.Value("admin").Boolean()
		user.NotContainsKey("password")
	}

	user := obj.Value(0).Object()
	user.Value("id").Number().IsEqual(1)
	user.Value("username").String().Contains("admin")
	user.Value("email").String().IsEqual("admin@example.com")
	user.Value("admin").Boolean().IsTrue()
	user.NotContainsKey("password")
}

func TestGetUsersWithFilters(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// search for admin users
	res := e.GET("/api/v1/users").
		WithQuery("isAdmin", "true").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)

	res.Header("X-Total-Count").AsNumber().IsEqual(1)

	obj := res.JSON().Array()
	obj.Length().IsEqual(1)

	user := obj.Value(0).Object()
	user.Value("id").Number().IsEqual(1)
	user.Value("username").String().Contains("admin")

	// search for the demo user
	res = e.GET("/api/v1/users").
		WithQuery("q", "demo").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)

	res.Header("X-Total-Count").AsNumber().IsEqual(1)

	obj = res.JSON().Array()
	obj.Length().IsEqual(1)

	user = obj.Value(0).Object()
	user.Value("id").Number().IsEqual(2)
	user.Value("username").String().Contains("demo")

	// search for a user that does not exist
	res = e.GET("/api/v1/users").
		WithQuery("q", "doesnotexist").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)
	res.Header("X-Total-Count").AsNumber().IsEqual(0)
}

func TestGetUser(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := e.GET("/api/v1/users/1").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	user.Value("id").Number().IsEqual(1)
	user.Value("username").String().Contains("admin")
	user.Value("email").String().IsEqual("admin@example.com")
	user.Value("admin").Boolean().IsTrue()
	user.NotContainsKey("password")
}

func TestCreateUpdateAndDeleteUser(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := e.POST("/api/v1/users").
		WithJSON(map[string]interface{}{
			"username": "Test User",
			"email":    "test@example.com",
		}).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	user.Value("id").Number().Gt(0)
	user.Value("username").String().IsEqual("test user")

	userId := user.Value("id").Number().Raw()
	userUrl := "/api/v1/users/" + strconv.FormatFloat(userId, 'f', -1, 64)

	user = e.GET(userUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	user.Value("id").Number().Gt(0)
	user.Value("username").String().IsEqual("test user")

	e.PUT(userUrl).
		WithJSON(map[string]interface{}{
			"username": "Test User Changed",
			"email":    "changed@example.com",
		}).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	user = e.GET(userUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	user.Value("id").Number().Gt(0)
	user.Value("username").String().IsEqual("test user changed")

	e.DELETE(userUrl).
		WithHeader("Authorization", "Bearer "+getJwtForAdminUser()).
		Expect().
		Status(http.StatusOK)

	e.GET(userUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusNotFound)
}

func TestUserAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", "/api/v1/users").Expect().Status(http.StatusUnauthorized)
	e.Request("GET", "/api/v1/users/1").Expect().Status(http.StatusUnauthorized)
	e.Request("POST", "/api/v1/users/").Expect().Status(http.StatusUnauthorized)
	e.Request("PUT", "/api/v1/users/1").Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", "/api/v1/users/1").Expect().Status(http.StatusUnauthorized)
}
