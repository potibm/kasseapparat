package tests_e2e

import (
	"net/http"
	"testing"
)

var purchaseStatsURL = "/api/v2/purchases/stats"

func TestPurchaseStats(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := e.GET(purchaseStatsURL).
		Expect().
		Status(http.StatusOK)

	res.Header("Access-Control-Allow-Origin").IsEqual("*")
	res.JSON().Object().Value("totalQuantity").Number()
	// store this value for later use
	totalQuantity := res.JSON().Object().Value("totalQuantity").Number().Raw()

	purchaseURL := createPurchase()

	res = e.GET(purchaseStatsURL).
		Expect().
		Status(http.StatusOK)

	res.JSON().Object().Value("totalQuantity").Number().IsEqual(totalQuantity + 1)

	deletePurchase(purchaseURL)
}
