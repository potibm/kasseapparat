package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	guestlistBaseUrl   = "/api/v2/guestlists"
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

	for i := range len(obj.Iter()) {
		list := obj.Value(i).Object()
		validateGuestlistObject(list)
	}

	list := obj.Value(0).Object()
	validateGuestlistObjectOne(list)
	list.Value("product").Object().Value("id").Number().IsEqual(2)
}

func TestGetGuestlistsEmpty(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(guestlistBaseUrl)).
		WithQuery("id", 0).
		Expect().
		Status(http.StatusOK)

	// assert that the response is an empty array
	res.JSON().Array().Length().IsEqual(0)
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

	res := withDemoUserAuthToken(e.GET(guestlistBaseUrl)).
		WithQuery("q", "Guestlist").
		Expect()

	res.Status(http.StatusOK)

	obj := res.JSON().Array()
	obj.NotEmpty()

	for _, item := range obj.Iter() {
		validateGuestlistObject(item.Object())
		item.Object().Value("name").String().Contains("Guestlist")
	}
}

func TestGetGuestlist(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	list := withDemoUserAuthToken(e.GET(guestlistUrlWithId)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	validateGuestlistObjectOne(list)
}

func TestCreateUpdateAndDeleteGuestList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	originalName := "Test List"

	changedName := "Test List Updated"

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
		Status(http.StatusNoContent)

	withDemoUserAuthToken(e.GET(listUrl)).
		Expect().
		Status(http.StatusNotFound)
}

func TestGuestlistAuthentication(t *testing.T) {
	testAuthenticationForEntityEndpoints(t, guestlistBaseUrl, guestlistUrlWithId)
}

func validateGuestlistObject(guestlist *httpexpect.Object) {
	guestlist.Value("id").Number().Gt(0)
	guestlist.Value("name").String().NotEmpty()
	guestlist.Value("typeCode").Boolean()
	guestlist.Value("productId").Number().Ge(0)
	guestlist.Value("product").Object()
}

func validateGuestlistObjectOne(guestlist *httpexpect.Object) {
	guestlist.Value("id").Number().IsEqual(1)
	guestlist.Value("name").String().Contains("Reduces")
	guestlist.Value("typeCode").Boolean().IsFalse()
	guestlist.Value("productId").Number().IsEqual(2)
}
