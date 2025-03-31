package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"
)

var purchaseBaseUrl = "/api/v2/purchases"

func TestGetPurchasesListIsEmpty(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get the purchase list
	purchaseListResponse := withDemoUserAuthToken(e.GET(purchaseBaseUrl)).
		Expect().
		Status(http.StatusOK)

	purchaseListResponse.Header(totalCountHeader).AsNumber().IsEqual(0)

	purchaseList := purchaseListResponse.JSON().Array()

	purchaseList.Length().IsEqual(0)
}

func TestGetPurchasesWithSort(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// define an array of sort fields
	sortFields := []string{"id", "createdAt", "totalGrossPrice", "createdBy.username"}

	for _, sortField := range sortFields {
		withDemoUserAuthToken(e.GET(purchaseBaseUrl)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestCreatePurchaseWithList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	purchaseResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalNetPrice":   "18.69",
			"totalGrossPrice": "20",
			"cart": []map[string]interface{}{
				{
					"ID":       2,
					"quantity": 1,
					"listItems": []map[string]interface{}{
						{
							"ID":             1,
							"attendedGuests": 1,
						},
					},
				},
			},
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	// we probably want to change this at one point :)
	purchaseResponse.Value("message").String().IsEqual("Purchase successful")
	purchaseResponse.Value("purchase").Object()
	purchase := purchaseResponse.Value("purchase").Object()
	purchase.Value("id").Number().Gt(0)
	purchase.Value("totalGrossPrice").String().IsEqual("20")
	purchase.Value("totalNetPrice").String().IsEqual("18.69")

	purchaseId := purchase.Value("id").Number().Raw()
	purchaseUrl := purchaseBaseUrl + "/" + strconv.FormatFloat(purchaseId, 'f', -1, 64)

	// Get the purchase
	purchase = withDemoUserAuthToken(e.GET(purchaseUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	purchase.Value("id").Number().IsEqual(purchaseId)
	purchase.Value("totalGrossPrice").String().IsEqual("20")
	purchase.Value("totalNetPrice").String().IsEqual("18.69")

	// Get the purchase list
	purchaseListResponse := withDemoUserAuthToken(e.GET(purchaseBaseUrl)).
		Expect().
		Status(http.StatusOK)

	purchaseListResponse.Header(totalCountHeader).AsNumber().IsEqual(1)

	purchaseList := purchaseListResponse.JSON().Array()

	purchaseList.Length().IsEqual(1)
	purchaseListItem := purchaseList.Value(0).Object()
	purchaseListItem.Value("id").Number().IsEqual(purchaseId)

	// Delete the purchase
	withDemoUserAuthToken(e.DELETE(purchaseUrl)).
		Expect().
		Status(http.StatusOK)

	withDemoUserAuthToken(e.GET(purchaseUrl)).
		Expect().
		Status(http.StatusNotFound)
}

func TestCreatePurchaseWithWrongTotalGrossPrice(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalGrossPrice": "21",
			"totalNetPrice":   "18.69",
			"cart": []map[string]interface{}{
				{
					"ID":        2,
					"quantity":  1,
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Total gross price does not match")
}

func TestCreatePurchaseWithWrongTotalNetPrice(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalGrossPrice": "20",
			"totalNetPrice":   "1.69",
			"cart": []map[string]interface{}{
				{
					"ID":        2,
					"quantity":  1,
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Total net price does not match")
}

func TestCreatePurchaseWithInvalidProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalGrossPrice": "21",
			"totalNetPrice":   "21",
			"cart": []map[string]interface{}{
				{
					"ID":        123,
					"quantity":  1,
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Product not found")
}

func TestCreatePurchaseWithListForWrongProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalNetPrice":   "11.11",
			"totalGrossPrice": "12.12",
			"cart": []map[string]interface{}{
				{
					"ID":       3,
					"quantity": 1,
					"listItems": []map[string]interface{}{
						{
							"ID":             1,
							"attendedGuests": 1,
						},
					},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "List item does not belong to product")
}

func TestCreatePurchaseWithListForAttendedGuestTooHigh(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalNetPrice":   "11.11",
			"totalGrossPrice": "12.12",
			"cart": []map[string]interface{}{
				{
					"ID":       3,
					"quantity": 1,
					"listItems": []map[string]interface{}{
						{
							"ID":             1,
							"attendedGuests": 15,
						},
					},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Additional guests exceed available guests")
}

func createPurchase() string {
	purchaseResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalNetPrice":   "37.38",
			"totalGrossPrice": "40",
			"cart": []map[string]interface{}{
				{
					"ID":        1,
					"quantity":  1,
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	purchase := purchaseResponse.Value("purchase").Object()
	purchaseId := purchase.Value("id").Number().Raw()

	return purchaseBaseUrl + "/" + strconv.FormatFloat(purchaseId, 'f', -1, 64)
}

func deletePurchase(purchaseUrl string) {
	withDemoUserAuthToken(e.DELETE(purchaseUrl)).
		Expect().
		Status(http.StatusOK)
}
