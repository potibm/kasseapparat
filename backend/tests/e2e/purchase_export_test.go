package tests_e2e

import (
	"net/http"
	"testing"
)

var purchaseExportBaseUrl = "/api/v2/purchases/export"

func TestPurchaseExportForEntityEndpoints(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	e.Request("GET", purchaseExportBaseUrl).Expect().Status(http.StatusUnauthorized)
}

func TestGetPurchaseExport(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	res := withDemoUserAuthToken(e.GET(purchaseExportBaseUrl)).
		Expect()

	res.Status(http.StatusOK)
	// .Contains("attachment; filename=purchase")
	// res.Header("Content-Disposition").Contains(".csv")
}
