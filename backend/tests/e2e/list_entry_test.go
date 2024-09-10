package tests_e2e

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestGetListEntries(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := e.GET("/api/v1/listEntries").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)

	totalCountHeader := res.Header("X-Total-Count").AsNumber()
	totalCountHeader.Gt(10)

	obj := res.JSON().Array()
	obj.Length().IsEqual(10)

	for i := 0; i < len(obj.Iter()); i++ {
		listEntry := obj.Value(i).Object()
		listEntry.Value("id").Number().Gt(0)
		listEntry.Value("name").String().NotEmpty()
		listEntry.Value("listId").Number().Gt(0)
		listEntry.Value("list").Object()
		code := listEntry.Value("code")
		if code.Raw() != nil {
			code.String().NotEmpty()
		} else {
			code.IsNull()
		}
		listEntry.Value("additionalGuests").Number().Ge(0)
		listEntry.Value("attendedGuests").Number().Ge(0)
		notifyOnArrivalEmail := listEntry.Value("notifyOnArrivalEmail")
		if notifyOnArrivalEmail.Raw() != nil {
			notifyOnArrivalEmail.String()
		} else {
			notifyOnArrivalEmail.IsNull()
		}
		purchaseId := listEntry.Value("purchaseId")
		if purchaseId.Raw() != nil {
			purchaseId.Number().Ge(0)
		} else {
			purchaseId.IsNull()
		}
	}

	listEntry := obj.Value(0).Object()
	listEntry.Value("id").Number().IsEqual(1)
	listEntry.Value("name").String().Length().Gt(5)
	listEntry.Value("code").IsNull()
	listEntry.Value("additionalGuests").Number().IsEqual(0)
	listEntry.Value("attendedGuests").Number().IsEqual(0)
	listEntry.Value("notifyOnArrivalEmail").IsNull()
	listEntry.Value("purchaseId").IsNull()
}

func TestGetListEntriesWithQuery(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := e.GET("/api/v1/listEntries").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect()

	res.Status(http.StatusOK)

	totalCountHeader := res.Header("X-Total-Count").AsNumber()
	totalCountHeader.Gt(10)

	obj := res.JSON().Array()

	listEntry := obj.Value(9).Object()
	name := listEntry.Value("name").String().Raw()

	name = strings.Split(name, " ")[0]
	log.Println("name: ", name)

	res = e.GET("/api/v1/listEntries").
		WithQuery("q", name).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK)

	res.JSON().Array()

	totalCountHeaderWithQuery := res.Header("X-Total-Count").AsNumber()
	totalCountHeaderWithQuery.Ge(1)
	totalCountHeaderWithQuery.Lt(totalCountHeader.Raw())

	obj = res.JSON().Array()
	obj.Length().Ge(1)

	for i := 0; i < len(obj.Iter()); i++ {
		listEntry := obj.Value(i).Object()
		listEntry.Value("name").String().Contains(name)
	}
}

func TestCreateUpdateAndDeleteListEntry(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	listEntry := e.POST("/api/v1/listEntries").
		WithJSON(map[string]interface{}{
			"listId":               1,
			"name":                 "Tessy Test",
			"additionalGuests":     2,
			"attendedGuests":       0,
			"arrivalNote":          "Hand out a tshirt",
			"notifyOnArrivalEmail": "test@example.com",
		}).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	listEntry.Value("id").Number().Gt(0)
	listEntry.Value("name").String().Contains("Tessy Test")
	listEntry.Value("additionalGuests").Number().IsEqual(2)
	listEntry.Value("attendedGuests").Number().IsEqual(0)
	listEntry.Value("arrivalNote").String().IsEqual("Hand out a tshirt")
	listEntry.Value("notifyOnArrivalEmail").String().IsEqual("test@example.com")

	listEntryId := listEntry.Value("id").Number().Raw()
	listEntryUrl := "/api/v1/listEntries/" + strconv.FormatFloat(listEntryId, 'f', -1, 64)

	listEntry = e.GET(listEntryUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	listEntry.Value("id").Number().Gt(0)
	listEntry.Value("name").String().Contains("Tessy Test")

	e.PUT(listEntryUrl).
		WithJSON(map[string]interface{}{
			"name":                 "Tessy Test Updated",
			"additionalGuests":     3,
			"attendedGuests":       0,
			"arrivalNote":          nil,
			"notifyOnArrivalEmail": nil,
		}).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	listEntry = e.GET(listEntryUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	listEntry.Value("id").Number().Gt(0)
	listEntry.Value("name").String().Contains("Tessy Test Updated")
	listEntry.Value("additionalGuests").Number().IsEqual(3)
	listEntry.Value("attendedGuests").Number().IsEqual(0)
	listEntry.Value("arrivalNote").IsNull()
	listEntry.Value("notifyOnArrivalEmail").IsNull()

	e.DELETE(listEntryUrl).
		WithHeader("Authorization", "Bearer "+getJwtForAdminUser()).
		Expect().
		Status(http.StatusOK)

	e.GET(listEntryUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusNotFound)

}

func TestListEntryAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", "/api/v1/listEntries").Expect().Status(http.StatusUnauthorized)
	e.Request("GET", "/api/v1/listEntries/1").Expect().Status(http.StatusUnauthorized)
	e.Request("POST", "/api/v1/listEntries/").Expect().Status(http.StatusUnauthorized)
	e.Request("PUT", "/api/v1/listEntries/1").Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", "/api/v1/listEntries/1").Expect().Status(http.StatusUnauthorized)
}
