package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	userBaseUrl   = "/api/v2/users"
	userUrlWithId = userBaseUrl + "/1"
)

func TestGetUsers(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(userBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Ge(2)

	obj := res.JSON().Array()
	obj.Length().Ge(2)

	for i := range len(obj.Iter()) {
		user := obj.Value(i).Object()
		validateUserObject(user)
	}

	user := obj.Value(0).Object()
	validateUserObjectAdmin(user)
}

func TestGetUsersEmpty(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(userBaseUrl)).
		WithQuery("q", "arandomstringthatnousernamecontains").
		Expect().
		Status(http.StatusOK)

	// assert that the response is an empty array
	res.JSON().Array().Length().IsEqual(0)
}

func TestGetUsersWithSort(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// define an array of sort fields
	sortFields := []string{"id", "username", "email", "admin"}

	for _, sortField := range sortFields {
		withDemoUserAuthToken(e.GET(userBaseUrl)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestGetUsersWithFilters(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// search for admin users
	res := withDemoUserAuthToken(e.GET(userBaseUrl)).
		WithQuery("isAdmin", "true").
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().IsEqual(1)

	obj := res.JSON().Array()
	obj.Length().IsEqual(1)

	user := obj.Value(0).Object()
	validateUserObjectAdmin(user)

	// search for the demo user
	res = withDemoUserAuthToken(e.GET(userBaseUrl)).
		WithQuery("q", "demo").
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().IsEqual(1)

	obj = res.JSON().Array()
	obj.Length().IsEqual(1)

	user = obj.Value(0).Object()
	validateUserObjectDemo(user)

	// search for a user that does not exist
	res = withDemoUserAuthToken(e.GET(userBaseUrl)).
		WithQuery("q", "doesnotexist").
		Expect()

	res.Status(http.StatusOK)
	res.Header(totalCountHeader).AsNumber().IsEqual(0)
}

func TestGetUser(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	user := withDemoUserAuthToken(e.GET(userUrlWithId)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	validateUserObjectAdmin(user)
}

func TestCreateUpdateAndDeleteUser(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	originalUsername := "test user"

	originalEmail := "test@example.com"

	changedUsername := "test user changed"

	changedEmail := "changed@example.com"

	user := withDemoUserAuthToken(e.POST(userBaseUrl)).
		WithJSON(map[string]interface{}{
			"username": originalUsername,
			"email":    originalEmail,
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	user.Value("id").Number().Gt(0)
	user.Value("username").String().IsEqual(originalUsername)
	user.Value("email").String().IsEqual(originalEmail)

	userId := user.Value("id").Number().Raw()
	userUrl := userBaseUrl + "/" + strconv.FormatFloat(userId, 'f', -1, 64)

	user = withDemoUserAuthToken(e.GET(userUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	user.Value("id").Number().Gt(0)
	user.Value("username").String().IsEqual(originalUsername)
	user.Value("email").String().IsEqual(originalEmail)

	withDemoUserAuthToken(e.PUT(userUrl)).
		WithJSON(map[string]interface{}{
			"username": changedUsername,
			"email":    changedEmail,
		}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	user = withDemoUserAuthToken(e.GET(userUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	user.Value("id").Number().Gt(0)
	user.Value("username").String().IsEqual(changedUsername)
	user.Value("email").String().IsEqual(changedEmail)

	withAdminUserAuthToken(e.DELETE(userUrl)).
		Expect().
		Status(http.StatusOK)

	withDemoUserAuthToken(e.GET(userUrl)).
		Expect().
		Status(http.StatusNotFound)
}

func TestUserAuthentication(t *testing.T) {
	testAuthenticationForEntityEndpoints(t, userBaseUrl, userUrlWithId)
}

func validateUserObject(user *httpexpect.Object) {
	user.Value("id").Number().Gt(0)
	user.Value("username").String().NotEmpty()
	user.Value("email").String().NotEmpty().Contains("@")
	user.Value("admin").Boolean()
	user.NotContainsKey("password")
}

func validateUserObjectAdmin(user *httpexpect.Object) {
	user.Value("id").Number().IsEqual(1)
	user.Value("username").String().Contains("admin")
	user.Value("email").String().IsEqual("admin@example.com")
	user.Value("admin").Boolean().IsTrue()
	user.NotContainsKey("password")
}

func validateUserObjectDemo(user *httpexpect.Object) {
	user.Value("id").Number().IsEqual(2)
	user.Value("username").String().Contains("demo")
	user.Value("email").String().IsEqual("demo@example.com")
	user.Value("admin").Boolean().IsFalse()
	user.NotContainsKey("password")
}
