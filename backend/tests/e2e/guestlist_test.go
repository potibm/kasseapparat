package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	guestlistBaseUrl   = "/api/v1/guestlists"
	guestlistUrlWithId = guestlistBaseUrl + "/1"
)

func TestGetGuestlists(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(guestlistBaseUrl)).
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

func TestGetGuestlistsWithSort(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// define an array of sort fields
	sortFields := []string{"id", "name"}

	for _, sortField := range sortFields {

		withDemoUserAuthToken(e.GET(guestlistBaseUrl)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestGetGuestlistsWithQuery(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	withDemoUserAuthToken(e.GET(guestlistBaseUrl)).
		WithQuery("q", "Guestlist").
		Expect().
		Status(http.StatusOK)

}

func TestGetList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	list := withDemoUserAuthToken(e.GET(guestlistUrlWithId)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	validateListObjectOne(list)
}

func TestCreateUpdateAndDelete(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	var originalName = "Test List"
	var changedName = "Test List Updated"

	list := withDemoUserAuthToken(e.POST(guestlistBaseUrl)).
		WithJSON(map[string]interface{}{
			"name":      originalName,
			"typeCode":  false,
			"productId": 2,
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	list.Value("id").Number().Gt(0)
	list.Value("name").String().IsEqual(originalName)

	listId := list.Value("id").Number().Raw()
	listUrl := guestlistBaseUrl + "/" + strconv.FormatFloat(listId, 'f', -1, 64)

	list = withDemoUserAuthToken(e.GET(listUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	list.Value("id").Number().Gt(0)
	list.Value("name").String().Contains(originalName)

	withDemoUserAuthToken(e.PUT(listUrl)).
		WithJSON(map[string]interface{}{
			"name":      changedName,
			"typeCode":  false,
			"productId": 2,
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
	testAuthenticationForEntityEndpoints(t, guestlistBaseUrl, guestlistUrlWithId)
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
