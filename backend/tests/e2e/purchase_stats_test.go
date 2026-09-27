package tests_e2e

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// publicEndpointToken guards GET /api/v3/purchases/stats, which stays reachable
// without a login. Kept here so init_test.go can put it in the test config.
const publicEndpointToken = "e2e-public-endpoint-token"

var purchaseStatsURL = "/api/v3/purchases/stats"

// withPublicToken authenticates a request against the token-guarded public endpoint,
// as the statistics display does: from another origin, with a bearer token.
func withPublicToken(req *httpexpect.Request) *httpexpect.Request {
	return req.
		WithHeader("Authorization", "Bearer "+publicEndpointToken).
		WithHeader("Origin", "https://some-other-site.example")
}

func TestPurchaseStatsRequiresToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// No token: the endpoint must not leak sales quantities to an anonymous caller.
	e.GET(purchaseStatsURL).
		Expect().
		Status(http.StatusUnauthorized)

	e.GET(purchaseStatsURL).
		WithHeader("Authorization", "Bearer wrong-token").
		Expect().
		Status(http.StatusUnauthorized)
}

func TestPurchaseStats(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withPublicToken(e.GET(purchaseStatsURL)).
		Expect().
		Status(http.StatusOK)

	res.Header("Access-Control-Allow-Origin").IsEqual("*")
	res.JSON().Object().Value("totalQuantity").Number()
	// store this value for later use
	totalQuantity := res.JSON().Object().Value("totalQuantity").Number().Raw()

	purchaseURL := createPurchase()

	res = withPublicToken(e.GET(purchaseStatsURL)).
		Expect().
		Status(http.StatusOK)

	res.JSON().Object().Value("totalQuantity").Number().IsEqual(totalQuantity + 1)

	deletePurchase(purchaseURL)
}
