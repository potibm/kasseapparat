package tests_e2e

import (
	"net/http"
	"testing"
)

var purchaseBaseURL = "/api/v2/purchases"

func TestGetPurchasesList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get the purchase list
	purchaseListResponse := withDemoUserAuthToken(e.GET(purchaseBaseURL)).
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

	purchaseListResponse := withDemoUserAuthToken(e.GET(purchaseBaseURL)).
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
		withDemoUserAuthToken(e.GET(purchaseBaseURL)).
			WithQuery("_sort", sortField).
			Expect().
			Status(http.StatusOK)
	}
}

func TestCreatePurchaseWithList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	purchaseResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CC",
			"totalNetPrice":   "18.69",
			"totalGrossPrice": "20",
			"cart": []map[string]any{
				{
					"ID":       2,
					"quantity": 1,
					"netPrice": "18.69",
					"listItems": []map[string]any{
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

	purchaseID := purchase.Value("id").String().Raw()
	purchaseURL := purchaseBaseURL + "/" + purchaseID

	// Get the purchase
	purchase = withDemoUserAuthToken(e.GET(purchaseURL)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	purchase.Value("id").String().IsEqual(purchaseID)
	purchase.Value("totalGrossPrice").String().IsEqual("20")
	purchase.Value("totalNetPrice").String().IsEqual("18.69")

	// Get the purchase list
	purchaseListResponse := withDemoUserAuthToken(e.GET(purchaseBaseURL)).
		WithQuery("_sort", "createdAt").
		WithQuery("_order", "DESC").
		Expect().
		Status(http.StatusOK)

	purchaseListResponse.Header(totalCountHeader).AsNumber().Ge(1)

	purchaseList := purchaseListResponse.JSON().Array()

	purchaseList.Length().Ge(1)
	purchaseListItem := purchaseList.Value(0).Object()
	purchaseListItem.Value("id").String().IsEqual(purchaseID)
	purchaseListItem.Value("paymentMethod").String().IsEqual("CC")

	// Delete the purchase
	withDemoUserAuthToken(e.DELETE(purchaseURL)).
		Expect().
		Status(http.StatusNoContent)

	withDemoUserAuthToken(e.GET(purchaseURL)).
		Expect().
		Status(http.StatusNotFound)
}

func TestCreatePurchaseWithWithFreeProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase with a free product
	purchaseResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CASH",
			"totalNetPrice":   "0",
			"totalGrossPrice": "0",
			"cart": []map[string]any{
				{
					"id":        3, // free product
					"quantity":  1,
					"netPrice":  "0",
					"listItems": []map[string]any{},
				},
			},
		}).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	purchase := purchaseResponse
	purchase.Value("id").String()
	purchase.Value("totalGrossPrice").String().IsEqual("0")
	purchase.Value("totalNetPrice").String().IsEqual("0")

	purchaseID := purchase.Value("id").String().Raw()
	purchaseURL := purchaseBaseURL + "/" + purchaseID

	// Get the purchase
	purchase = withDemoUserAuthToken(e.GET(purchaseURL)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	purchase.Value("id").String().IsEqual(purchaseID)
	purchase.Value("totalGrossPrice").String().IsEqual("0")
	purchase.Value("totalNetPrice").String().IsEqual("0")

	deletePurchase(purchaseURL)
}

func TestCreatePurchaseWithWrongTotalGrossPrice(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total gross price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CASH",
			"totalGrossPrice": "21",
			"totalNetPrice":   "18.69",
			"cart": []map[string]any{
				{
					"ID":        2,
					"quantity":  1,
					"netPrice":  "18.69",
					"listItems": []map[string]any{},
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
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CASH",
			"totalGrossPrice": "20",
			"totalNetPrice":   "1.69",
			"cart": []map[string]any{
				{
					"ID":        2,
					"quantity":  1,
					"netPrice":  "18.69",
					"listItems": []map[string]any{},
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
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CASH",
			"totalGrossPrice": "20",
			"totalNetPrice":   "18.69",
			"cart": []map[string]any{
				{
					"ID":        2,
					"quantity":  1,
					"netPrice":  "1.69", // wrong price
					"listItems": []map[string]any{},
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
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CASH",
			"totalGrossPrice": "21",
			"totalNetPrice":   "21",
			"cart": []map[string]any{
				{
					"ID":        123,
					"quantity":  1,
					"netPrice":  "21",
					"listItems": []map[string]any{},
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
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "INVALID",
			"totalGrossPrice": "21",
			"totalNetPrice":   "21",
			"cart": []map[string]any{
				{
					"ID":        123,
					"quantity":  1,
					"netPrice":  "21",
					"listItems": []map[string]any{},
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
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CASH",
			"totalNetPrice":   "0",
			"totalGrossPrice": "0",
			"cart": []map[string]any{
				{
					"ID":       3, // free product
					"quantity": 1,
					"netPrice": 0,
					"listItems": []map[string]any{
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
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CASH",
			"totalNetPrice":   "0",
			"totalGrossPrice": "0",
			"cart": []map[string]any{
				{
					"ID":       3, // free product
					"quantity": 1,
					"netPrice": 0,
					"listItems": []map[string]any{
						{
							"ID":             1,
							"attendedGuests": 10,
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

	purchaseURLWithID := createPurchase()

	e.Request("GET", purchaseBaseURL).Expect().Status(http.StatusUnauthorized)
	e.Request("GET", purchaseURLWithID).Expect().Status(http.StatusUnauthorized)
	e.Request("POST", purchaseBaseURL).Expect().Status(http.StatusUnauthorized)
	e.Request("DELETE", purchaseURLWithID).Expect().Status(http.StatusUnauthorized)

	deletePurchase(purchaseURLWithID)
}

func TestPurchaseGetByIdWithoutUuid(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get a purchase with a wrong ID
	errorResponse := withDemoUserAuthToken(e.GET(purchaseBaseURL + "/123")).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Invalid purchase ID")
}

func TestPurchaseDeleteWithoutUuid(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get a purchase with a wrong ID
	errorResponse := withDemoUserAuthToken(e.DELETE(purchaseBaseURL + "/123")).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	validateErrorDetailMessage(errorResponse, "Invalid purchase ID")
}

func createPurchase() string {
	purchaseResponse := withDemoUserAuthToken(e.POST(purchaseBaseURL)).
		WithJSON(map[string]any{
			"paymentMethod":   "CASH",
			"totalNetPrice":   "37.38",
			"totalGrossPrice": "40",
			"cart": []map[string]any{
				{
					"ID":        1,
					"quantity":  1,
					"netPrice":  "37.38",
					"listItems": []map[string]any{},
				},
			},
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	purchase := purchaseResponse
	purchaseID := purchase.Value("id").String().Raw()

	return purchaseBaseURL + "/" + purchaseID
}

func deletePurchase(purchaseURL string) {
	withDemoUserAuthToken(e.DELETE(purchaseURL)).
		Expect().
		Status(http.StatusNoContent)
}
