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
	guestsImportUrl       = "/api/v1/guestsUpload"
	guestsImportCsvHeader = "Code;LastName;FirstName;Subject;Blocked\n"
	deineTicketProductId  = 4
)

func TestGuestImport(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	fileContent := guestsImportCsvHeader +
		"123456789;XYZTEST Lastname;Firstname;EV123;\n"

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

	withDemoUserAuthToken(e.POST(guestsImportUrl)).
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
		"12345678A;XYZTEST Lastname;Firstname;EV123;\n"

	guestImportResponse := uploadGuestImport(fileContent).
		Status(http.StatusOK).
		JSON().
		Object()

	guestImportResponse.Value("createdGuests").Number().IsEqual(1)
	guestImportResponse.Value("warnings").Array().IsEmpty()

	fileContent = guestsImportCsvHeader +
		"12345678A;XYZTEST Dupe;Dupesten;EV123;\n" +
		"123;XYZTEST Invalid;Codesten;EV123;\n" +
		"ABCABCABC;XYZTEST Blocked;Dupesten;EV123;blocked\n"

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
	guestsUrl := productBaseUrl + "/" + strconv.Itoa(deineTicketProductId) + "/guests"
	guests := withDemoUserAuthToken(e.GET(guestsUrl)).
		WithQuery("q", query).
		Expect().
		Status(http.StatusOK).JSON().Array()

	for i := 0; i < len(guests.Iter()); i++ {
		guest := guests.Value(i).Object()
		guestId := guest.Value("id").Number().Raw()
		log.Println("Deleting list entry with id", guestId)

		withDemoUserAuthToken(e.DELETE(guestBaseUrl + "/" + strconv.Itoa(int(guestId)))).
			Expect().
			Status(http.StatusOK)
	}
}

func uploadGuestImport(fileContent string) *httpexpect.Response {
	reader := strings.NewReader(fileContent)

	// Create a list entry import
	return withDemoUserAuthToken(e.POST(guestsImportUrl)).
		WithMultipart().
		WithFile("file", "import.csv", reader).
		Expect()
}
