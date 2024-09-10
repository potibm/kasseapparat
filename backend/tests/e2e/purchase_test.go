package tests_e2e

import (
	"net/http"
	"strconv"
	"testing"
)

func TestCreatePurchaseWithList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	purchaseResponse := e.POST("/api/v1/purchases").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
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
	purchaseUrl := "/api/v1/purchases/" + strconv.FormatFloat(purchaseId, 'f', -1, 64)

	// Get the purchase
	purchase = e.GET(purchaseUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK).JSON().Object()

	purchase.Value("id").Number().IsEqual(purchaseId)
	purchase.Value("totalPrice").Number().IsEqual(20.0)

	// Delete the purchase
	e.DELETE(purchaseUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusOK)

	e.GET(purchaseUrl).
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
		Expect().
		Status(http.StatusNotFound)
}

func TestCreatePurchaseWithWrongTotalPrice(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total price
	errrorResponse := e.POST("/api/v1/purchases").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
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

	// we probably want to change this at one point :)
	errrorResponse.Value("details").String().IsEqual("Total price does not match")
}

func TestCreatePurchaseWithInvalidProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase, but with a wrong total price
	errrorResponse := e.POST("/api/v1/purchases").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
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

	errrorResponse.Value("details").String().IsEqual("Product not found")
}

func TestCreatePurchaseWithListForWrongProduct(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	errrorResponse := e.POST("/api/v1/purchases").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
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

	errrorResponse.Value("details").String().IsEqual("List item does not belong to product")
}

func TestCreatePurchaseWithListForAttendedGuestTooHigh(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a purchase
	errrorResponse := e.POST("/api/v1/purchases").
		WithHeader("Authorization", "Bearer "+getJwtForDemoUser()).
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

	errrorResponse.Value("details").String().IsEqual("Additional guests exceed available guests")
}
