package tests_e2e

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	listEntryBaseUrl   = "/api/v1/listEntries"
	listEntryUrlWithId = listEntryBaseUrl + "/1"
)

func TestGetListEntries(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(listEntryBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Gt(10)

	obj := res.JSON().Array()
	obj.Length().IsEqual(10)

	for i := 0; i < len(obj.Iter()); i++ {
		listEntry := obj.Value(i).Object()
		validateListEntryObject(listEntry)
	}

	listEntry := obj.Value(0).Object()
	validateListEntryObjectOne(listEntry)
}

func TestGetListEntriesWithQuery(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(listEntryBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	totalCountHeaderValue := res.Header(totalCountHeader).AsNumber()
	totalCountHeaderValue.Gt(10)

	obj := res.JSON().Array()

	listEntry := obj.Value(9).Object()
	name := listEntry.Value("name").String().Raw()

	name = strings.Split(name, " ")[0]
	log.Println("name: ", name)

	res = withDemoUserAuthToken(e.GET(listEntryBaseUrl)).
		WithQuery("q", name).
		Expect().
		Status(http.StatusOK)

	res.JSON().Array()

	totalCountHeaderWithQuery := res.Header(totalCountHeader).AsNumber()
	totalCountHeaderWithQuery.Ge(1)
	totalCountHeaderWithQuery.Lt(totalCountHeaderValue.Raw())

	obj = res.JSON().Array()
	obj.Length().Ge(1)

	for i := 0; i < len(obj.Iter()); i++ {
		listEntry := obj.Value(i).Object()
		listEntry.Value("name").String().Contains(name)
	}
}

func TestGetListEntry(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	listEntry := withDemoUserAuthToken(e.GET(listEntryUrlWithId)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	validateListEntryObjectOne(listEntry)
}

func TestCreateUpdateAndDeleteListEntry(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	originalName := "Tessy Test"
	changedName := "Tessy Test Updated"
	notifyEmail := "test@example.com"
	arrivalNote := "Hand out a tshirt"

	listEntry := withDemoUserAuthToken(e.POST("/api/v1/listEntries")).
		WithJSON(map[string]interface{}{
			"listId":               1,
			"name":                 originalName,
			"additionalGuests":     2,
			"attendedGuests":       0,
			"arrivalNote":          arrivalNote,
			"notifyOnArrivalEmail": notifyEmail,
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	listEntry.Value("id").Number().Gt(0)
	listEntry.Value("name").String().Contains(originalName)
	listEntry.Value("additionalGuests").Number().IsEqual(2)
	listEntry.Value("attendedGuests").Number().IsEqual(0)
	listEntry.Value("arrivalNote").String().IsEqual(arrivalNote)
	listEntry.Value("notifyOnArrivalEmail").String().IsEqual(notifyEmail)

	listEntryId := listEntry.Value("id").Number().Raw()
	listEntryUrl := listEntryBaseUrl + "/" + strconv.FormatFloat(listEntryId, 'f', -1, 64)

	listEntry = withDemoUserAuthToken(e.GET(listEntryUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	listEntry.Value("id").Number().Gt(0)
	listEntry.Value("name").String().Contains(originalName)

	withDemoUserAuthToken(e.PUT(listEntryUrl)).
		WithJSON(map[string]interface{}{
			"name":                 changedName,
			"additionalGuests":     3,
			"attendedGuests":       0,
			"arrivalNote":          nil,
			"notifyOnArrivalEmail": nil,
		}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	listEntry = withDemoUserAuthToken(e.GET(listEntryUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	listEntry.Value("id").Number().Gt(0)
	listEntry.Value("name").String().Contains(changedName)
	listEntry.Value("additionalGuests").Number().IsEqual(3)
	listEntry.Value("attendedGuests").Number().IsEqual(0)
	listEntry.Value("arrivalNote").IsNull()
	listEntry.Value("notifyOnArrivalEmail").IsNull()

	withDemoUserAuthToken(e.DELETE(listEntryUrl)).
		Expect().
		Status(http.StatusOK)

	withDemoUserAuthToken(e.GET(listEntryUrl)).
		Expect().
		Status(http.StatusNotFound)

}

func TestListEntryAuthentication(t *testing.T) {
	testAuthenticationForEntityEndpoints(t, listEntryBaseUrl, listEntryUrlWithId)
}

func TestGetListEntriesByProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	url := "/api/v1/products/2/listEntries"

	res := withDemoUserAuthToken(e.GET(url)).
		WithQuery("q", "e").
		Expect().
		Status(http.StatusOK)

	obj := res.JSON().Array()
	obj.Length().Ge(1)

	for i := 0; i < len(obj.Iter()); i++ {
		listEntry := obj.Value(i).Object()
		validateListEntrySummaryObject(listEntry)

		name := listEntry.Value("name").String().Raw()
		if !strings.Contains(strings.ToLower(name), "e") {
			t.Errorf("Name does not contain 'e' or 'E': %s", name)
		}
	}
}

func validateListEntryBasicObject(listEntry *httpexpect.Object) {
	listEntry.Value("id").Number().Gt(0)
	listEntry.Value("name").String().NotEmpty()
	code := listEntry.Value("code")
	if code.Raw() != nil {
		code.String().NotEmpty()
	} else {
		code.IsNull()
	}
	listEntry.Value("additionalGuests").Number().Ge(0)
	arrivalNote := listEntry.Value("arrivalNote")
	if arrivalNote.Raw() != nil {
		arrivalNote.String()
	} else {
		arrivalNote.IsNull()
	}
}

func validateListEntrySummaryObject(listEntry *httpexpect.Object) {
	validateListEntryBasicObject(listEntry)
	listEntry.Value("listName").String().NotEmpty()
}

func validateListEntryObject(listEntry *httpexpect.Object) {
	validateListEntryBasicObject(listEntry)
	listEntry.Value("listId").Number().Gt(0)
	listEntry.Value("list").Object()
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

func validateListEntryObjectOne(listEntry *httpexpect.Object) {
	listEntry.Value("id").Number().IsEqual(1)
	listEntry.Value("name").String().Length().Gt(5)
	listEntry.Value("code").IsNull()
	listEntry.Value("additionalGuests").Number().IsEqual(0)
	listEntry.Value("attendedGuests").Number().IsEqual(0)
	listEntry.Value("notifyOnArrivalEmail").IsNull()
	listEntry.Value("purchaseId").IsNull()
}
