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
	guestBaseUrl   = "/api/v2/guests"
	guestUrlWithId = guestBaseUrl + "/1"
)

func TestGetGuests(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(guestBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().Gt(10)

	obj := res.JSON().Array()
	obj.Length().IsEqual(10)

	for i := range len(obj.Iter()) {
		guest := obj.Value(i).Object()
		validateGuestObject(guest)
	}

	guest := obj.Value(0).Object()
	validateGuestObjectOne(guest)
}

func TestGetGuestWithQuery(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(guestBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	totalCountHeaderValue := res.Header(totalCountHeader).AsNumber()
	totalCountHeaderValue.Gt(10)

	obj := res.JSON().Array()

	guest := obj.Value(9).Object()
	name := guest.Value("name").String().Raw()

	name = strings.Split(name, " ")[0]
	log.Println("name: ", name)

	res = withDemoUserAuthToken(e.GET(guestBaseUrl)).
		WithQuery("q", name).
		Expect().
		Status(http.StatusOK)

	res.JSON().Array()

	totalCountHeaderWithQuery := res.Header(totalCountHeader).AsNumber()
	totalCountHeaderWithQuery.Ge(1)
	totalCountHeaderWithQuery.Lt(totalCountHeaderValue.Raw())

	obj = res.JSON().Array()
	obj.Length().Ge(1)

	for i := range len(obj.Iter()) {
		guest := obj.Value(i).Object()
		guest.Value("name").String().ContainsFold(name)
	}
}

func TestGetGuestWithSort(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// define an array of sort fields
	sortFields := []string{"id", "name", "guestlist.name", "arrivedAt"}

	for _, sortField := range sortFields {
		withDemoUserAuthToken(e.GET(guestBaseUrl)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestGetGuestsWithQueryIsPresent(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// isPresent
	res := withDemoUserAuthToken(e.GET(guestBaseUrl)).
		WithQuery("isPresent", "true").
		Expect()

	res.Status(http.StatusOK)

	obj := res.JSON().Array()
	obj.NotEmpty()

	for _, item := range obj.Iter() {
		validateGuestObject(item.Object())
		item.Object().Value("arrivedAt").String().NotEmpty()
	}
}

func TestGetGuestsWithQueryIsNotPresent(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// isPresent
	res := withDemoUserAuthToken(e.GET(guestBaseUrl)).
		WithQuery("isNotPresent", "true").
		Expect()

	res.Status(http.StatusOK)

	obj := res.JSON().Array()
	obj.NotEmpty()

	for _, item := range obj.Iter() {
		validateGuestObject(item.Object())
		item.Object().Value("arrivedAt").IsNull()
	}
}

func TestGetGuestsWithQueryGuestlist(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// isPresent
	res := withDemoUserAuthToken(e.GET(guestBaseUrl)).
		WithQuery("guestlist_id", 1).
		Expect()

	res.Status(http.StatusOK)

	obj := res.JSON().Array()
	obj.NotEmpty()

	for _, item := range obj.Iter() {
		validateGuestObject(item.Object())
		item.Object().Value("guestlistId").IsEqual(1)
		item.Object().Value("guestlist").Object().Value("id").IsEqual(1)
	}
}

func TestGetGuest(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	guest := withDemoUserAuthToken(e.GET(guestUrlWithId)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	validateGuestObjectOne(guest)
}

func TestCreateUpdateAndDeleteGuest(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	originalName := "Tessy Test"
	changedName := "Tessy Test Updated"
	notifyEmail := "test@example.com"
	arrivalNote := "Hand out a tshirt"

	guest := withDemoUserAuthToken(e.POST(guestBaseUrl)).
		WithJSON(map[string]interface{}{
			"guestlistId":          1,
			"name":                 originalName,
			"additionalGuests":     2,
			"attendedGuests":       0,
			"arrivalNote":          arrivalNote,
			"notifyOnArrivalEmail": notifyEmail,
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	guest.Value("id").Number().Gt(0)
	guest.Value("name").String().Contains(originalName)
	guest.Value("additionalGuests").Number().IsEqual(2)
	guest.Value("attendedGuests").Number().IsEqual(0)
	guest.Value("arrivalNote").String().IsEqual(arrivalNote)
	guest.Value("notifyOnArrivalEmail").String().IsEqual(notifyEmail)

	guestId := guest.Value("id").Number().Raw()
	guestUrl := guestBaseUrl + "/" + strconv.FormatFloat(guestId, 'f', -1, 64)

	guest = withDemoUserAuthToken(e.GET(guestUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	guest.Value("id").Number().Gt(0)
	guest.Value("name").String().Contains(originalName)

	withDemoUserAuthToken(e.PUT(guestUrl)).
		WithJSON(map[string]interface{}{
			"name":                 changedName,
			"additionalGuests":     3,
			"attendedGuests":       0,
			"arrivalNote":          nil,
			"notifyOnArrivalEmail": nil,
		}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	guest = withDemoUserAuthToken(e.GET(guestUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	guest.Value("id").Number().Gt(0)
	guest.Value("name").String().Contains(changedName)
	guest.Value("additionalGuests").Number().IsEqual(3)
	guest.Value("attendedGuests").Number().IsEqual(0)
	guest.Value("arrivalNote").IsNull()
	guest.Value("notifyOnArrivalEmail").IsNull()

	withDemoUserAuthToken(e.DELETE(guestUrl)).
		Expect().
		Status(http.StatusOK)

	withDemoUserAuthToken(e.GET(guestUrl)).
		Expect().
		Status(http.StatusNotFound)
}

func TestGuestAuthentication(t *testing.T) {
	testAuthenticationForEntityEndpoints(t, guestBaseUrl, guestUrlWithId)
}

func TestGuestsByProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	url := productBaseUrl + "/2/guests"

	res := withDemoUserAuthToken(e.GET(url)).
		WithQuery("q", "e").
		Expect().
		Status(http.StatusOK)

	obj := res.JSON().Array()
	obj.Length().Ge(1)

	for i := range len(obj.Iter()) {
		guest := obj.Value(i).Object()
		validateGuestSummaryObject(guest)

		name := guest.Value("name").String().Raw()
		if !strings.Contains(strings.ToLower(name), "e") {
			t.Errorf("Name does not contain 'e' or 'E': %s", name)
		}
	}
}

func validateGuestBasicObject(guest *httpexpect.Object) {
	guest.Value("id").Number().Gt(0)
	guest.Value("name").String().NotEmpty()

	code := guest.Value("code")
	if code.Raw() != nil {
		code.String().NotEmpty()
	} else {
		code.IsNull()
	}

	guest.Value("additionalGuests").Number().Ge(0)

	arrivalNote := guest.Value("arrivalNote")
	if arrivalNote.Raw() != nil {
		arrivalNote.String()
	} else {
		arrivalNote.IsNull()
	}
}

func validateGuestSummaryObject(guest *httpexpect.Object) {
	validateGuestBasicObject(guest)
	guest.Value("listName").String().NotEmpty()
}

func validateGuestObject(guest *httpexpect.Object) {
	validateGuestBasicObject(guest)
	guest.Value("guestlistId").Number().Gt(0)
	guest.Value("guestlist").Object()
	guest.Value("attendedGuests").Number().Ge(0)

	notifyOnArrivalEmail := guest.Value("notifyOnArrivalEmail")
	if notifyOnArrivalEmail.Raw() != nil {
		notifyOnArrivalEmail.String()
	} else {
		notifyOnArrivalEmail.IsNull()
	}

	purchaseId := guest.Value("purchaseId")
	if purchaseId.Raw() != nil {
		purchaseId.Number().Ge(0)
	} else {
		purchaseId.IsNull()
	}
}

func validateGuestObjectOne(guest *httpexpect.Object) {
	guest.Value("id").Number().IsEqual(1)
	guest.Value("name").String().Length().Gt(5)
	guest.Value("code").IsNull()
	guest.Value("additionalGuests").Number().Ge(0)
	guest.Value("attendedGuests").Number().IsEqual(0)
	guest.Value("notifyOnArrivalEmail").IsNull()
	guest.Value("purchaseId").IsNull()
}
