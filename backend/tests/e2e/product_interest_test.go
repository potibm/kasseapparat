package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	productInterestBaseUrl   = "/api/v1/productInterests"
	productInterestUrlWithId = productInterestBaseUrl + "/1"
)

func TestGetProductInterest(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(productInterestBaseUrl)).
		Expect()

	res.Status(http.StatusOK)

	res.Header(totalCountHeader).AsNumber().IsEqual(0)

	res.JSON().Array().IsEmpty()
}

func TestGetProductInterestsWithSort(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// define an array of sort fields
	sortFields := []string{"id", "pos", "createdAt", "product.id", "product.name"}

	for _, sortField := range sortFields {

		withDemoUserAuthToken(e.GET(productInterestBaseUrl)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestCreateAndDeleteProductInterest(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	productInterest := withDemoUserAuthToken(e.POST(productInterestBaseUrl)).
		WithJSON(map[string]interface{}{
			"productId": 1,
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	productInterest.Value("id").Number().Gt(0)

	productInterestId := productInterest.Value("id").Number().Raw()
	productInterestUrl := productInterestBaseUrl + "/" + strconv.FormatFloat(productInterestId, 'f', -1, 64)

	getTotalCountOfProductInterests().IsEqual(1)

	withDemoUserAuthToken(e.DELETE(productInterestUrl)).
		Expect().
		Status(http.StatusOK)

	getTotalCountOfProductInterests().IsEqual(0)

}

func getTotalCountOfProductInterests() *httpexpect.Number {
	res := withDemoUserAuthToken(e.GET(productInterestBaseUrl)).
		Expect().
		Status(http.StatusOK)

	return res.Header(totalCountHeader).AsNumber()
}

func TestProductInterestAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", productInterestBaseUrl).Expect().Status(http.StatusUnauthorized)
	e.Request("POST", productInterestBaseUrl).Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", productInterestUrlWithId).Expect().Status(http.StatusUnauthorized)
}
