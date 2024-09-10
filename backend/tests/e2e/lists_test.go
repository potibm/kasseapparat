package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"
)

func TestGetLists(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := e.GET("/api/v1/lists").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)

	totalCountHeader := res.Header("X-Total-Count").AsNumber()
	totalCountHeader.Gt(10)

	obj := res.JSON().Array()
	obj.Length().IsEqual(10)

	for i := 0; i < len(obj.Iter()); i++ {
		list := obj.Value(i).Object()
		list.Value("id").Number().Gt(0)
		list.Value("name").String().NotEmpty()
		list.Value("typeCode").Boolean()
		list.Value("productId").Number().Ge(0)
		list.Value("product").Object()
	}

	list := obj.Value(0).Object()
	list.Value("id").Number().IsEqual(1)
	list.Value("name").String().Contains("Reduces")
	list.Value("typeCode").Boolean().IsFalse()
	list.Value("productId").Number().IsEqual(2)
	list.Value("product").Object().Value("id").Number().IsEqual(2)
}

func TestGetList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	list := e.GET("/api/v1/lists/1").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	list.Value("id").Number().IsEqual(1)
	list.Value("name").String().Contains("Reduces")
	list.Value("typeCode").Boolean().IsFalse()
	list.Value("productId").Number().IsEqual(2)
}

func TestCreateUpdateAndDeleteList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	list := e.POST("/api/v1/lists").
		WithJSON(map[string]interface{}{
			"name":      "Test List",
			"TypeCode":  false,
			"ProductId": 2,
		}).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	list.Value("id").Number().Gt(0)
	list.Value("name").String().Contains("Test List")

	listId := list.Value("id").Number().Raw()
	listUrl := "/api/v1/lists/" + strconv.FormatFloat(listId, 'f', -1, 64)

	list = e.GET(listUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	list.Value("id").Number().Gt(0)
	list.Value("name").String().Contains("Test List")

	e.PUT(listUrl).
		WithJSON(map[string]interface{}{
			"name":      "Test List Updated",
			"TypeCode":  false,
			"ProductId": 2,
		}).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	list = e.GET(listUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	list.Value("id").Number().Gt(0)
	list.Value("name").String().Contains("Test List Updated")

	e.DELETE(listUrl).
		WithHeader("Authorization", "Bearer "+getJwtForAdminUser()).
		Expect().
		Status(http.StatusOK)

	e.GET(listUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusNotFound)

}

func TestListAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", "/api/v1/lists").Expect().Status(http.StatusUnauthorized)
	e.Request("GET", "/api/v1/lists/1").Expect().Status(http.StatusUnauthorized)
	e.Request("POST", "/api/v1/lists/").Expect().Status(http.StatusUnauthorized)
	e.Request("PUT", "/api/v1/lists/1").Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", "/api/v1/lists/1").Expect().Status(http.StatusUnauthorized)
}
