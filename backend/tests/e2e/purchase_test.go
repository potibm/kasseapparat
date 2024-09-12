package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"
)

var (
	purchaseBaseUrl = "/api/v1/purchases"
)

func TestCreatePurchaseWithList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	purchaseResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalPrice": 20.0,
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
	purchase.Value("totalPrice").Number().IsEqual(20.0)

	purchaseId := purchase.Value("id").Number().Raw()
	purchaseUrl := purchaseBaseUrl + "/" + strconv.FormatFloat(purchaseId, 'f', -1, 64)

	// Get the purchase
	purchase = withDemoUserAuthToken(e.GET(purchaseUrl)).
		Expect().
		Status(http.StatusOK).JSON().Object()

	purchase.Value("id").Number().IsEqual(purchaseId)
	purchase.Value("totalPrice").Number().IsEqual(20.0)

	// Delete the purchase
	withDemoUserAuthToken(e.DELETE(purchaseUrl)).
		Expect().
		Status(http.StatusOK)

	withDemoUserAuthToken(e.GET(purchaseUrl)).
		Expect().
		Status(http.StatusNotFound)
}

func TestCreatePurchaseWithWrongTotalPrice(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalPrice": 21.0,
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

	validateErrorDetailMessage(errorResponse, "Total price does not match")
}

func TestCreatePurchaseWithInvalidProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total price
	errorResponse := withDemoUserAuthToken(e.POST(purchaseBaseUrl)).
		WithJSON(map[string]interface{}{
			"totalPrice": 21.0,
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
			"totalPrice": 0.0,
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
			"totalPrice": 0.0,
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
			"totalPrice": 40.0,
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
