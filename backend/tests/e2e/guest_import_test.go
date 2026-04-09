package tests_e2e

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	guestsImportURL       = "/api/v2/guestsUpload"
	guestsImportCsvHeader = "Code;LastName;FirstName;Subject;Blocked;Notiz;\n"
	deineTicketProductID  = 4
)

func TestGuestImport(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	fileContent := guestsImportCsvHeader +
		"123456789;XYZTEST Lastname;Firstname;EV123;;T-shirt size XL;\n"

	guestImportResponse := uploadGuestImport(fileContent).
		Status(http.StatusOK).
		JSON().
		Object()

	guestImportResponse.Value("createdGuests").Number().IsEqual(1)
	guestImportResponse.Value("warnings").Array().IsEmpty()

	deleteGuestsByNameQuery("XYZTEST")
}

func TestGuestImportWithoutFile(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	withDemoUserAuthToken(e.POST(guestsImportURL)).
		Expect().
		Status(http.StatusBadRequest)
}

func TestGuestImportWithEmptyFile(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	fileContent := ""

	guestImportResponse := uploadGuestImport(fileContent).
		Status(http.StatusBadRequest).
		JSON().
		Object()

	guestImportResponse.Value("details").String().Contains("Failed to read header")
}

func TestGuestImportWithWarningMessages(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	fileContent := guestsImportCsvHeader +
		"12345678A;XYZTEST Lastname;Firstname;EV123;;;\n"

	guestImportResponse := uploadGuestImport(fileContent).
		Status(http.StatusOK).
		JSON().
		Object()

	guestImportResponse.Value("createdGuests").Number().IsEqual(1)
	guestImportResponse.Value("warnings").Array().IsEmpty()

	fileContent = guestsImportCsvHeader +
		"12345678A;XYZTEST Dupe;Dupesten;EV123;;;\n" +
		"123;XYZTEST Invalid;Codesten;EV123;;T-shirt size M;\n" +
		"ABCABCABC;XYZTEST Blocked;Dupesten;EV123;blocked;;\n"

	guestImportResponse = uploadGuestImport(fileContent).
		Status(http.StatusOK).
		JSON().
		Object()

	guestImportResponse.Value("createdGuests").Number().IsEqual(0)
	guestImportResponse.Value("warnings").Array().Length().IsEqual(3)
	guestImportWarnings := guestImportResponse.Value("warnings").Array()

	guestImportWarnings.Value(0).String().IsEqual("Already exists: 12345678A (1)")
	guestImportWarnings.Value(1).String().IsEqual("Invalid code: 123 (2)")
	guestImportWarnings.Value(2).String().IsEqual("Blocked: ABCABCABC (3)")

	deleteGuestsByNameQuery("XYZTEST")
}

func deleteGuestsByNameQuery(query string) {
	guestsURL := productBaseURL + "/" + strconv.Itoa(deineTicketProductID) + "/guests"
	guests := withDemoUserAuthToken(e.GET(guestsURL)).
		WithQuery("q", query).
		Expect().
		Status(http.StatusOK).JSON().Array()

	for i := range len(guests.Iter()) {
		guest := guests.Value(i).Object()
		guestID := guest.Value("id").Number().Raw()

		withDemoUserAuthToken(e.DELETE(guestBaseURL + "/" + strconv.Itoa(int(guestID)))).
			Expect().
			Status(http.StatusOK)
	}
}

func uploadGuestImport(fileContent string) *httpexpect.Response {
	reader := strings.NewReader(fileContent)

	// Create a list entry import
	return withDemoUserAuthToken(e.POST(guestsImportURL)).
		WithMultipart().
		WithFile("file", "import.csv", reader).
		Expect()
}

func TestGuestsImportAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("POST", guestsImportURL).Expect().Status(http.StatusUnauthorized)
}
