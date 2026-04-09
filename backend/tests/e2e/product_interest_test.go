package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var (
	productInterestBaseURL   = "/api/v2/productInterests"
	productInterestURLWithID = productInterestBaseURL + "/1"
)

func TestGetProductInterest(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(productInterestBaseURL)).
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
		withDemoUserAuthToken(e.GET(productInterestBaseURL)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestCreateAndDeleteProductInterest(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	productInterest := withDemoUserAuthToken(e.POST(productInterestBaseURL)).
		WithJSON(map[string]any{
			"productId": 1,
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	productInterest.Value("id").Number().Gt(0)

	productInterestID := productInterest.Value("id").Number().Raw()
	productInterestURL := productInterestBaseURL + "/" + strconv.FormatFloat(productInterestID, 'f', -1, 64)

	getTotalCountOfProductInterests().IsEqual(1)

	withDemoUserAuthToken(e.DELETE(productInterestURL)).
		Expect().
		Status(http.StatusNoContent)

	getTotalCountOfProductInterests().IsEqual(0)
}

func getTotalCountOfProductInterests() *httpexpect.Number {
	res := withDemoUserAuthToken(e.GET(productInterestBaseURL)).
		Expect().
		Status(http.StatusOK)

	return res.Header(totalCountHeader).AsNumber()
}

func TestProductInterestAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", productInterestBaseURL).Expect().Status(http.StatusUnauthorized)
	e.Request("POST", productInterestBaseURL).Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", productInterestURLWithID).Expect().Status(http.StatusUnauthorized)
}
