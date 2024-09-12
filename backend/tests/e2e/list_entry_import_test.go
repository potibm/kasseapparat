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
	listEntryImportUrl       = "/api/v1/listEntriesUpload"
	listEntryImportCsvHeader = "Code;LastName;FirstName;Subject;Blocked\n"
	deineTicketProductId     = 4
)

func TestListEntryImport(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	fileContent := listEntryImportCsvHeader +
		"123456789;XYZTEST Lastname;Firstname;EV123;\n"

	listEntryImportResponse := uploadListEntryImport(fileContent).
		Status(http.StatusOK).
		JSON().
		Object()

	listEntryImportResponse.Value("createdEntries").Number().IsEqual(1)
	listEntryImportResponse.Value("warnings").Array().IsEmpty()

	deleteListEntriesByNameQuery("XYZTEST")
}

func TestListEntryImportWithoutFile(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	withDemoUserAuthToken(e.POST(listEntryImportUrl)).
		Expect().
		Status(http.StatusBadRequest)
}

func TestListEntryImportWithEmptyFile(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	fileContent := ""

	listEntryImportResponse := uploadListEntryImport(fileContent).
		Status(http.StatusBadRequest).
		JSON().
		Object()

	listEntryImportResponse.Value("details").String().Contains("Failed to read header")
}

func TestListEntryImportWithWarningMessages(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	fileContent := listEntryImportCsvHeader +
		"12345678A;XYZTEST Lastname;Firstname;EV123;\n"

	listEntryImportResponse := uploadListEntryImport(fileContent).
		Status(http.StatusOK).
		JSON().
		Object()

	listEntryImportResponse.Value("createdEntries").Number().IsEqual(1)
	listEntryImportResponse.Value("warnings").Array().IsEmpty()

	fileContent = listEntryImportCsvHeader +
		"12345678A;XYZTEST Dupe;Dupesten;EV123;\n" +
		"123;XYZTEST Invalid;Codesten;EV123;\n" +
		"ABCABCABC;XYZTEST Blocked;Dupesten;EV123;blocked\n"

	listEntryImportResponse = uploadListEntryImport(fileContent).
		Status(http.StatusOK).
		JSON().
		Object()

	listEntryImportResponse.Value("createdEntries").Number().IsEqual(0)
	listEntryImportResponse.Value("warnings").Array().Length().IsEqual(3)
	listEntryImportWarnings := listEntryImportResponse.Value("warnings").Array()

	listEntryImportWarnings.Value(0).String().IsEqual("Already exists: 12345678A (1)")
	listEntryImportWarnings.Value(1).String().IsEqual("Invalid code: 123 (2)")
	listEntryImportWarnings.Value(2).String().IsEqual("Blocked: ABCABCABC (3)")

	deleteListEntriesByNameQuery("XYZTEST")

}

func deleteListEntriesByNameQuery(query string) {
	listEntryUrl := productBaseUrl + "/" + strconv.Itoa(deineTicketProductId) + "/listEntries"
	listEntries := withDemoUserAuthToken(e.GET(listEntryUrl)).
		WithQuery("q", query).
		Expect().
		Status(http.StatusOK).JSON().Array()

	for i := 0; i < len(listEntries.Iter()); i++ {
		listEntry := listEntries.Value(i).Object()
		listEntryId := listEntry.Value("id").Number().Raw()
		log.Println("Deleting list entry with id", listEntryId)

		withDemoUserAuthToken(e.DELETE("/api/v1/listEntries/" + strconv.Itoa(int(listEntryId)))).
			Expect().
			Status(http.StatusOK)
	}
}

func uploadListEntryImport(fileContent string) *httpexpect.Response {
	reader := strings.NewReader(fileContent)

	// Create a list entry import
	return withDemoUserAuthToken(e.POST(listEntryImportUrl)).
		WithMultipart().
		WithFile("file", "import.csv", reader).
		Expect()
}
