package tests_e2e

import (
	"net/http"
	"testing"
)

var (
	sumupReadersURL       = "/api/v3/sumup/readers"
	sumupReadersURLWithID = sumupReadersURL + "/reader_1"
)

func TestGetSumupReaders(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(sumupReadersURL)).
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().IsEqual(1)

	obj := res.JSON().Array()
	obj.Length().IsEqual(1)

	for i := range len(obj.Iter()) {
		reader := obj.Value(i).Object()
		reader.Value("id").String().NotEmpty()
		reader.Value("name").String().IsEqual("Mock Reader 1")
	}
}

func TestGetSumupReader(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(sumupReadersURLWithID)).
		Expect()

	res.Status(http.StatusOK)

	reader := res.JSON().Object()
	reader.Value("id").String().IsEqual("reader_1")
	reader.Value("name").String().IsEqual("Mock Reader")
}

func TestDeleteSumupReader(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	withDemoUserAuthToken(e.DELETE(sumupReadersURLWithID)).
		Expect().
		Status(http.StatusNoContent)
}

// Note: Authentication tests removed for Phase 1 - auth is now handled by reverse proxy in Phase 2
