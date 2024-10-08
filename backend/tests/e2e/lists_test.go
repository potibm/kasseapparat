package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	listBaseUrl   = "/api/v1/lists"
	listUrlWithId = listBaseUrl + "/1"
)

func TestGetLists(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(listBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Gt(10)

	obj := res.JSON().Array()
	obj.Length().IsEqual(10)

	for i := 0; i < len(obj.Iter()); i++ {
		list := obj.Value(i).Object()
		validateListObject(list)
	}

	list := obj.Value(0).Object()
	validateListObjectOne(list)
	list.Value("product").Object().Value("id").Number().IsEqual(2)
}

func TestGetListsWithSort(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// define an array of sort fields
	sortFields := []string{"id", "name"}

	for _, sortField := range sortFields {

		withDemoUserAuthToken(e.GET(listBaseUrl)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestGetList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	list := withDemoUserAuthToken(e.GET(listUrlWithId)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	validateListObjectOne(list)
}

func TestCreateUpdateAndDeleteList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	var originalName = "Test List"
	var changedName = "Test List Updated"

	list := withDemoUserAuthToken(e.POST(listBaseUrl)).
		WithJSON(map[string]interface{}{
			"name":      originalName,
			"TypeCode":  false,
			"ProductId": 2,
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	list.Value("id").Number().Gt(0)
	list.Value("name").String().IsEqual(originalName)

	listId := list.Value("id").Number().Raw()
	listUrl := listBaseUrl + "/" + strconv.FormatFloat(listId, 'f', -1, 64)

	list = withDemoUserAuthToken(e.GET(listUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	list.Value("id").Number().Gt(0)
	list.Value("name").String().Contains(originalName)

	withDemoUserAuthToken(e.PUT(listUrl)).
		WithJSON(map[string]interface{}{
			"name":      changedName,
			"TypeCode":  false,
			"ProductId": 2,
		}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	list = withDemoUserAuthToken(e.GET(listUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	list.Value("id").Number().Gt(0)
	list.Value("name").String().Contains(changedName)

	withDemoUserAuthToken(e.DELETE(listUrl)).
		Expect().
		Status(http.StatusOK)

	withDemoUserAuthToken(e.GET(listUrl)).
		Expect().
		Status(http.StatusNotFound)
}

func TestListAuthentication(t *testing.T) {
	testAuthenticationForEntityEndpoints(t, listBaseUrl, listUrlWithId)
}

func validateListObject(list *httpexpect.Object) {
	list.Value("id").Number().Gt(0)
	list.Value("name").String().NotEmpty()
	list.Value("typeCode").Boolean()
	list.Value("productId").Number().Ge(0)
	list.Value("product").Object()
}

func validateListObjectOne(list *httpexpect.Object) {
	list.Value("id").Number().IsEqual(1)
	list.Value("name").String().Contains("Reduces")
	list.Value("typeCode").Boolean().IsFalse()
	list.Value("productId").Number().IsEqual(2)
}
