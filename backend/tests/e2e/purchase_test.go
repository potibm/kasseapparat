package tests_e2e

import (
	"net/http"
	"testing"
)

var purchaseBaseUrl = "/api/v2/purchases"

// @TODO: test a purchase with only a free product, currently leads to an error

func TestGetPurchasesList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get the purchase list
	purchaseListResponse := withDemoUserAuthToken(e.GET(purchaseBaseUrl)).
		Expect().
		Status(http.StatusOK)

	purchaseListResponse.Header(totalCountHeader).AsNumber().Ge(1)

	purchaseList := purchaseListResponse.JSON().Array()

	purchaseList.Length().Ge(1)

	purchaseListItem := purchaseList.Value(0).Object()
	purchaseListItem.Value("id").String()
	purchaseListItem.Value("totalGrossPrice").String()
	purchaseListItem.Value("totalNetPrice").String()
	purchaseListItem.Value("totalVatAmount").String()
	purchaseListItem.Value("createdAt").String().NotEmpty()
	purchaseListItem.Value("createdBy").Object().Value("username").String().NotEmpty()
	purchaseListItem.Value("paymentMethod").String().NotEmpty()
	purchaseListItem.Value("purchaseItems").Array().Length().Gt(0)
}

func TestGetPurchasesListWithAllFilters(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	purchaseListResponse := withDemoUserAuthToken(e.GET(purchaseBaseUrl)).
		WithQuery("paymentMethod", "CASH").
		WithQuery("createdById", "1").
		WithQuery("totalGrossPrice_gte", "1").
		WithQuery("totalGrossPrice_lte", "100").
		WithQuery("id", "1,2,3").
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
	sortFields := []string{"id", "createdAt", "totalGrossPrice", "createdBy.username", "paymentMethod"}

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
			"paymentMethod":   "CC",
			"totalNetPrice":   "18.69",
			"totalGrossPrice": "20",
			"cart": []map[string]interface{}{
				{
					"ID":       2,
					"quantity": 1,
					"netPrice": "18.69",
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

	purchase := purchaseResponse
	purchase.Value("id").String()
	purchase.Value("totalGrossPrice").String().IsEqual("20")
	purchase.Value("totalNetPrice").String().IsEqual("18.69")
	purchaseItems := purchase.Value("purchaseItems").Array()
	purchaseItems.Length().IsEqual(1)
	purchaseItem := purchaseItems.Value(0).Object()
	purchaseItem.Value("id").Number().Gt(0)
	purchaseItem.Value("productID").Number().IsEqual(2)
	product := purchaseItem.Value("product").Object()
	product.Value("id").Number().IsEqual(2)
	product.Value("name").String().NotEmpty() // see issue #469

	purchaseId := purchase.Value("id").String().Raw()
	purchaseUrl := purchaseBaseUrl + "/" + purchaseId

	// Get the purchase
	purchase = withDemoUserAuthToken(e.GET(purchaseUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	purchase.Value("id").String().IsEqual(purchaseId)
	purchase.Value("totalGrossPrice").String().IsEqual("20")
	purchase.Value("totalNetPrice").String().IsEqual("18.69")

	// Get the purchase list
	purchaseListResponse := withDemoUserAuthToken(e.GET(purchaseBaseUrl)).
		WithQuery("_sort", "createdAt").
		WithQuery("_order", "DESC").
		Expect().
		Status(http.StatusOK)

	purchaseListResponse.Header(totalCountHeader).AsNumber().Ge(1)

	purchaseList := purchaseListResponse.JSON().Array()

	purchaseList.Length().Ge(1)
	purchaseListItem := purchaseList.Value(0).Object()
	purchaseListItem.Value("id").String().IsEqual(purchaseId)
	purchaseListItem.Value("paymentMethod").String().IsEqual("CC")

	// Delete the purchase
	withDemoUserAuthToken(e.DELETE(purchaseUrl)).
		Expect().
		Status(http.StatusNoContent)

	withDemoUserAuthToken(e.GET(purchaseUrl)).
		Expect().
		Status(http.StatusNotFound)
}

func TestCreatePurchaseWithWrongTotalGrossPrice(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total gross price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"paymentMethod":   "CASH",
			"totalGrossPrice": "21",
			"totalNetPrice":   "18.69",
			"cart": []map[string]interface{}{
				{
					"ID":        2,
					"quantity":  1,
					"netPrice":  "18.69",
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

	// Create a purchase, but with a wrong total net price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"paymentMethod":   "CASH",
			"totalGrossPrice": "20",
			"totalNetPrice":   "1.69",
			"cart": []map[string]interface{}{
				{
					"ID":        2,
					"quantity":  1,
					"netPrice":  "18.69",
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Total net price does not match")
}

func TestCreatePurchaseWithWrongProductPrice(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total net price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"paymentMethod":   "CASH",
			"totalGrossPrice": "20",
			"totalNetPrice":   "18.69",
			"cart": []map[string]interface{}{
				{
					"ID":        2,
					"quantity":  1,
					"netPrice":  "1.69", // wrong price
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Invalid product price")
}

func TestCreatePurchaseWithInvalidProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong product ID
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"paymentMethod":   "CASH",
			"totalGrossPrice": "21",
			"totalNetPrice":   "21",
			"cart": []map[string]interface{}{
				{
					"ID":        123,
					"quantity":  1,
					"netPrice":  "21",
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Product not found")
}

func TestCreatePurchaseWithInvalidPaymentMethod(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong payment method
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"paymentMethod":   "INVALID",
			"totalGrossPrice": "21",
			"totalNetPrice":   "21",
			"cart": []map[string]interface{}{
				{
					"ID":        123,
					"quantity":  1,
					"netPrice":  "21",
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "invalid payment method")
}

func TestCreatePurchaseWithListForWrongProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"paymentMethod":   "CASH",
			"totalNetPrice":   "0",
			"totalGrossPrice": "0",
			"cart": []map[string]interface{}{
				{
					"ID":       3, // free product
					"quantity": 1,
					"netPrice": 0,
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
			"paymentMethod":   "CASH",
			"totalNetPrice":   "0",
			"totalGrossPrice": "0",
			"cart": []map[string]interface{}{
				{
					"ID":       3, // free product
					"quantity": 1,
					"netPrice": 0,
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

func TestPurchasesAuthentication(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	purchaseUrlWithId := createPurchase()

	e.Request("GET", purchaseBaseUrl).Expect().Status(http.StatusUnauthorized)
	e.Request("GET", purchaseUrlWithId).Expect().Status(http.StatusUnauthorized)
	e.Request("POST", purchaseBaseUrl).Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", purchaseUrlWithId).Expect().Status(http.StatusUnauthorized)

	deletePurchase(purchaseUrlWithId)
}

func TestPurchaseGetByIdWithoutUuid(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get a purchase with a wrong ID
	errorResponse := withDemoUserAuthToken(e.GET(purchaseBaseUrl + "/123")).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Invalid purchase ID")
}

func TestPurchaseDeleteWithoutUuid(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get a purchase with a wrong ID
	errorResponse := withDemoUserAuthToken(e.DELETE(purchaseBaseUrl + "/123")).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Invalid purchase ID")
}

func createPurchase() string {
	purchaseResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"paymentMethod":   "CASH",
			"totalNetPrice":   "37.38",
			"totalGrossPrice": "40",
			"cart": []map[string]interface{}{
				{
					"ID":        1,
					"quantity":  1,
					"netPrice":  "37.38",
					"listItems": []map[string]interface{}{},
				},
			},
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	purchase := purchaseResponse
	purchaseId := purchase.Value("id").String().Raw()

	return purchaseBaseUrl + "/" + purchaseId
}

func deletePurchase(purchaseUrl string) {
	withDemoUserAuthToken(e.DELETE(purchaseUrl)).
		Expect().
		Status(http.StatusNoContent)
}
