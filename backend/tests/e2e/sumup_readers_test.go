package tests_e2e

import (
	"net/http"
	"testing"
)

var (
	sumupReadersUrl       = "/api/v2/sumup/readers"
	sumupReadersUrlWithId = sumupReadersUrl + "/reader_1"
)

func TestGetSumupReaders(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(sumupReadersUrl)).
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

	res := withDemoUserAuthToken(e.GET(sumupReadersUrlWithId)).
		Expect()

	res.Status(http.StatusOK)

	reader := res.JSON().Object()
	reader.Value("id").String().IsEqual("reader_1")
	reader.Value("name").String().IsEqual("Mock Reader")
}

func TestDeleteSumupReader(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	withDemoUserAuthToken(e.DELETE(sumupReadersUrlWithId)).
		Expect().
		Status(http.StatusNoContent)
}

func TestSumupReadersAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", sumupReadersUrl).Expect().Status(http.StatusUnauthorized)
	e.Request("POST", sumupReadersUrl).Expect().Status(http.StatusUnauthorized)
	e.Request("GET", sumupReadersUrlWithId).Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", sumupReadersUrlWithId).Expect().Status(http.StatusUnauthorized)
}
